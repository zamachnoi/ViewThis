package data

import (
	"errors"

	"github.com/zamachnoi/viewthis/lib"
	"github.com/zamachnoi/viewthis/models"
	"gorm.io/gorm"
)

func GetAllQueues() ([]models.Queue, error) {
    var queues []models.Queue
    err := lib.GetDB().Find(&queues).Error
    if err != nil {
        return nil, err
    }
    if len(queues) == 0 {
        return nil, errors.New("no queues found")
    }
    return queues, nil
}

func CreateQueue(queue models.Queue) (*models.Queue, error) {
	if err := lib.GetDB().Create(&queue).Error; err != nil {
		return nil, err
	}
	return &queue, nil
}

// GetQueue returns the queue with the given ID and all submissions in it
func GetQueueByID(id uint) (*models.Queue, error) {
    var queue models.Queue
    if err := lib.GetDB().Preload("Submissions").First(&queue, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &queue, nil
}

func GetQueueByName(name string) (*models.Queue, error) {
    var queue models.Queue
    err := lib.GetDB().Where("name LIKE ?", name).First(&queue).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &queue, nil
}

func ClearQueueByID(id uint) error {
    var queue models.Queue
    if err := lib.GetDB().First(&queue, id).Error; err != nil {
        return err
    }
    if err := lib.GetDB().Model(&queue).Association("Submissions").Clear(); err != nil {
        return err
    }
    return nil
}

func UpdateQueue(queue models.Queue) (*models.Queue, error) {
    if err := lib.GetDB().Save(&queue).Error; err != nil {
        return nil, err
    }
    return &queue, nil
}

func DeleteQueue(id uint) error {
	var queue models.Queue
    if err := lib.GetDB().First(&queue, id).Error; err != nil {
        return err
    }
    if err := lib.GetDB().Delete(&queue).Error; err != nil {
        return err
    }
    return nil

}