package clihandler

const (
	takeOrderCommandText   = "takeorder"
	returnOrderCommandText = "returnorder"
	giveOrderCommandText   = "giveorder"
	orderListCommandText   = "orderlist"
	limitArgText           = "limit"
	instockArgText         = "instock"
	takeRefundCommandText  = "takerefund"
	refundListCommandText  = "refundlist"
	checkPvzCommandText    = "checkpvz"
	newPvzCommandText      = "newpvz"
)

// HelpCommandText is the name for help command
const HelpCommandText = "help"

const (
	helpIDpvz = "%IDпвз"
	helpText  = "Команды:\n" +
		helpIDpvz + " " + takeOrderCommandText + " %IDзаказа %IDклиента %Дата(в форме YYYY-MM-DD) %ВесЗаказа %ЦенаЗаказа %Упаковка(опционально)\n" +
		" - Принять заказ от курьера;\n" +
		helpIDpvz + " " + returnOrderCommandText + " %IDзаказа\n" +
		" - Вернуть заказ курьеру;\n" +
		helpIDpvz + " " + giveOrderCommandText + " %IDзаказа1 %IDзаказа2 %IDзаказа3...\n" +
		" - Выдать заказ(ы) клиенту;\n" +
		helpIDpvz + " " + orderListCommandText + " %IDклиента " + limitArgText + " %N " + instockArgText + "\n" +
		" - Получить список заказов, limit %N и instock опциональны\n" +
		" - limit %N выводит N последних записей вместо всех\n" +
		" - instock выводит заказы находящиеся в ПВЗ\n" +
		helpIDpvz + " " + takeRefundCommandText + " %IDклиента %IDзаказа\n" +
		" - Принять возврат от клиента;\n" +
		helpIDpvz + " " + refundListCommandText + " %страница %размерстраницы\n" +
		" - Получить страницу из списка возвратов.\n" +
		helpIDpvz + " " + checkPvzCommandText +
		" - Вывести информацию о ПВЗ;\n" +
		helpIDpvz + " " + newPvzCommandText + " \"%имя\" \"%адрес\" \"%контакты\""
)
