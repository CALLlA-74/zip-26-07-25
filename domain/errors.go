package domain

import "errors"

var (
	//ErrInternalServerError = errors.New("Внутренняя ошибка сервера")
	ErrBusyServer       = errors.New("В данный момент сервер занят")
	ErrAddingImpossible = errors.New("Архив обрабатывается -- добавление новых файлов невозможно")
	ErrTaskNotFound     = errors.New("Задача не найдена")
)
