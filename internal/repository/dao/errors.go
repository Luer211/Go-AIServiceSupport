package dao

import (
	"errors"
	"gorm.io/gorm"
	"fmt"
)


var (
	ErrNotFound		  = errors.New("dao: record not found")
	ErrAlreadyExists  = errors.New("dao: record already exists")
	ErrNotImplemented = errors.New("dao: method not implemented")
)


func wrapDBError(operation string, err error) error {
    // 把原生错误，统一转换成业务层可识别的标准化错误
	
	if err == nil {
        return nil
    }

    switch {
	// gorm层错误，数据不存在
    case errors.Is(err, gorm.ErrRecordNotFound):
        return fmt.Errorf("%s: %w", operation, ErrNotFound)

	// gorm层错误，数据已存在
    case errors.Is(err, gorm.ErrDuplicatedKey):
        return fmt.Errorf("%s: %w", operation, ErrAlreadyExists)

	// 其它底层错误，直接暴露
    default:
        return fmt.Errorf("%s: %w", operation, err)
    }
}