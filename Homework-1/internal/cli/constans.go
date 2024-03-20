package cli

const (
	takeOrderCommandText   = "takeorder"
	returnOrderCommandText = "returnorder"
	giveOrderCommandText   = "giveorder"
	orderListCommandText   = "orderlist"
	takeRefundCommandText  = "takerefund"
	refundListCommandText  = "refundlist"
)

const HelpCommandText = "help"

const helpText = "Команды:\n" +
	takeOrderCommandText + " %IDзаказа %IDклиента %Дата(в форме YYYY-MM-DD)\n" +
	" - Принять заказ от курьера;\n" +
	returnOrderCommandText + " %IDзаказа\n" +
	" - Вернуть заказ курьеру;\n" +
	giveOrderCommandText + " %IDзаказа1 %IDзаказа2 %IDзаказа3...\n" +
	" - Выдать заказ(ы) клиенту;\n" +
	orderListCommandText + " %IDклиента limit %N instock\n" +
	" - Получить список заказов, limit %N и instock опциональны\n" +
	" - limit %N выводит N последних записей вместо всех\n" +
	" - instock выводит заказы находящиеся в ПВЗ\n" +
	takeRefundCommandText + " %IDклиента %IDзаказа\n" +
	" - Принять возврат от клиента;\n" +
	refundListCommandText + " %страница %размерстраницы\n" +
	" - Получить страницу из списка возвратов."
