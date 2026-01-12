package service

import (
	"fmt"
	"sum11crawler/models"
	"sync"

	"gorm.io/gorm"
)

type SumService struct {
	db     *gorm.DB
	parser *SumParser
}

func NewSumService(db *gorm.DB, parser *SumParser) *SumService {
	return &SumService{db, parser}
}

func (s *SumService) ParseAndStoreArticles(messagesF func(string), percentageF func(models.ParsingPercentage)) {
	var total int64
	s.db.Model(&models.Link{}).Where("(html IS NOT NULL AND html != ?) AND type = ?", "", LinkTypeArticle).Count(&total)

	var items []models.Link
	q := s.db.Where("(html IS NOT NULL AND html != ?) AND desc IS NULL AND type = ?", "", LinkTypeArticle)
	err := q.Find(&items).Error
	if err != nil {
		panic(err)
	}

	itemsCount := len(items)
	doneCount := total - int64(itemsCount)

	messagesF(fmt.Sprintf("Selected articles: %d", itemsCount))

	partCount := int(doneCount)

	goroutineCount := 10
	if itemsCount < goroutineCount {
		goroutineCount = itemsCount
	}
	if goroutineCount <= 0 {
		messagesF("Nothing to parse")
		return
	}

	var parsedItems []models.Link
	var parsingMu sync.Mutex
	var parsingWG sync.WaitGroup

	itemsPerGoroutine := itemsCount / goroutineCount

	err = s.db.Transaction(func(tx *gorm.DB) error {
		for n := range goroutineCount {
			from := n * itemsPerGoroutine
			to := (n + 1) * itemsPerGoroutine

			var goroutineItems []models.Link
			if n+1 == goroutineCount {
				goroutineItems = items[from:]
			} else {
				goroutineItems = items[from:to]
			}

			parsingWG.Add(1)

			go func() {
				defer parsingWG.Done()
				for _, item := range goroutineItems {
					if item.HTML == nil {
						messagesF("No HTML")
						continue
					}
					word, title, desc, err := s.parser.ParseArticle(*item.HTML)
					if err != nil {
						messagesF(err.Error())
					}

					item.Word = &word
					item.Desc = &desc
					item.Title = &title

					parsingMu.Lock()
					tx.Save(item)
					// err = tx.Save(item).Error
					// if err != nil {
					// 	return err
					// }
					parsedItems = append(parsedItems, item)
					partCount++
					parsingMu.Unlock()

					percentage := s.calcPercentage(partCount, int(total))

					// Output percentage data
					percentageF(models.ParsingPercentage{
						Percentage: percentage,
						Word:       word,
						URL:        item.URL,
					})
				}
			}()
		}

		parsingWG.Wait()
		messagesF("Writing...")

		return nil
	})
}

func (s *SumService) calcPercentage(part, total int) float32 {
	return float32(part) / float32(total) * 100
}
