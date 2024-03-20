package cli

const (
	checkPvzCommandText    = "checkpvz"
	newPvzCommandText      = "newpvz"
	takeOrderCommandText   = "takeorder"
	returnOrderCommandText = "returnorder"
	giveOrderCommandText   = "giveorder"
	orderListCommandText   = "orderlist"
	takeRefundCommandText  = "takerefund"
	refundListCommandText  = "refundlist"
	//exitCommandText        = "exit"
)

const HelpCommandText = "help"

const helpText = "Команды:\n" +
	"%IDпвз " + checkPvzCommandText +
	" - Вывести информацию о ПВЗ;\n" +
	"%IDпвз " + newPvzCommandText + " %имя / %адрес / %контакты\n" +
	" - Добавить ПВЗ в базу данных;\n" +
	" - Разделите информацию знаками / и пробелами;\n" +
	"%IDпвз " + takeOrderCommandText + " %IDзаказа %IDклиента %Дата(в форме YYYY-MM-DD)\n" +
	" - Принять заказ от курьера;\n" +
	"%IDпвз " + returnOrderCommandText + " %IDзаказа\n" +
	" - Вернуть заказ курьеру;\n" +
	"%IDпвз " + giveOrderCommandText + " %IDзаказа1 %IDзаказа2 %IDзаказа3...\n" +
	" - Выдать заказ(ы) клиенту;\n" +
	"%IDпвз " + orderListCommandText + " %IDклиента limit %N instock\n" +
	" - Получить список заказов, limit %N и instock опциональны\n" +
	" - limit %N выводит N последних записей вместо всех\n" +
	" - instock выводит заказы находящиеся в ПВЗ\n" +
	"%IDпвз " + takeRefundCommandText + " %IDклиента %IDзаказа\n" +
	" - Принять возврат от клиента;\n" +
	"%IDпвз " + refundListCommandText + " %страница %размерстраницы\n" +
	" - Получить страницу из списка возвратов."

const pvzStorageFile = "../data/pvzinfo.json"
