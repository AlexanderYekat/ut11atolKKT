﻿#Область ВстариваниеВКонфуУТ11
функция ПолчитьИмяОбработкиДляПечатиЧека()
	//"C:\docum\fl\utAtol\НастройкаСвязисККТ_fl.epf"
	//Сообщить("Процедура ПолчитьИмяОбработкиДляПечатиЧека");									
	возврат "C:\docum\fl\utAtol\ПечатьЧекаККТВнеКонфигурации.epf";
конецФункции //ПолчитьИмяОбработкиДляПечатиЧека

Функция ПечатьВызватьКомандуНаККТВместоВстроеннойВКонфигурацию(Команда, ОповещениеПриЗавершенииПечатиЧека, ФормаДокумента,
																	докСсылкаПоКоторомуПечатаемЧек = "", ИмяКассира = "") экспорт
	//"C:\docum\fl\utAtol\НастройкаСвязисККТ_fl.epf"									
	//Сообщить("Процедура ПечатьЧекаНаККТВместоВстроеннойВКонфигурацию начало");
	имяОбработкиДляПечатиЧека = ПолчитьИмяОбработкиДляПечатиЧека();
	допПараметры = Новый Структура;
	допПараметры.Вставить("Команда", Команда);
	допПараметры.Вставить("ФормаДокум", ФормаДокумента);
	допПараметры.Вставить("Кассир", ИмяКассира);
	//проуедура обрабатывающая результат печти чека - она находится в модуле формы документа
	допПараметры.Вставить("ОповещениеПриЗавершенииПечатиЧека", ОповещениеПриЗавершенииПечатиЧека);
	//данные необходимые для формирования чека
	допПараметры.Вставить("докСсылкаПоКоторомуПечатаемЧек", докСсылкаПоКоторомуПечатаемЧек);
    ОбработкаОкончанияПомещения = Новый ОписаниеОповещения
        ("ОбработчикОкончанияПомещения", ЭтотОбъект, допПараметры);    
    НачатьПомещениеФайла(ОбработкаОкончанияПомещения, , 
        имяОбработкиДляПечатиЧека, Ложь, ФормаДокумента.УникальныйИдентификатор);	
	//Сообщить("ПечатьЧекаНаККТВместоВстроеннойВКонфигурацию конец");
	возврат ИСТИНА;
	//если вернём ложь, то продолжится стандартная обработка печати чека конфигурации
КонецФункции //ПечатьЧекаНаККТВместоВстроеннойВКонфигурацию ()

Процедура ОбработчикОкончанияПомещения(Результат, Адрес, 
        ВыбранноеИмяФайла, ДополнительныеПараметры) экспорт
	//Сообщить("Процедура ОбработчикОкончанияПомещения");
	Если НЕ Результат Тогда
		//Сообщить("Процедура ОбработчикОкончанияПомещения --");
        Сообщить("Не удалось открыть обработку печати чека");
		возврат;
	КонецЕсли; 
	имяОбр = МенеджерОборудованияВызовСервера.ПолучитьИмяВнешнейОбработки(Адрес);	
	ПарамДляОбработкиПечатиЧека = Новый Структура;
	//передаём данные на основе которых можно было бы сформировать чека
	ПарамДляОбработкиПечатиЧека.Вставить("докСсылкаПоКоторомуПечатаемЧек", ДополнительныеПараметры.докСсылкаПоКоторомуПечатаемЧек);
	ПарамДляОбработкиПечатиЧека.Вставить("Команда", ДополнительныеПараметры.Команда);
	ПарамДляОбработкиПечатиЧека.Вставить("Кассир", ДополнительныеПараметры.Кассир);
	////передаём форму документа, что спразу же во внешней обработке её модифицировать (записать номер чека)
	//ПарамДляОбработкиПечатиЧека.Вставить("ФормаДокум", ДополнительныеПараметры.ФормаДокум);
	//при закрытии формы документа, форма печати чека тоже закроется
	ФормаОжиданияОтвета = ПолучитьФорму("ВнешняяОбработка."+имяОбр+".Форма.Форма", ПарамДляОбработкиПечатиЧека, ДополнительныеПараметры.ФормаДокум);
	//после закрытия формы печати чека - вызывает процедуру обработки результата печати чека
	ФормаОжиданияОтвета.ОписаниеОповещенияОЗакрытии = ДополнительныеПараметры.ОповещениеПриЗавершенииПечатиЧека;
	//открываем форму, которая нам напечатает чек
	ФормаОжиданияОтвета.Открыть();	
	//Сообщить("Процедура ОбработчикОкончанияПомещения --");
КонецПроцедуры //ОбработчикОкончанияПомещения

#КонецОбласти   

Процедура ОбработчикОкончанияПомещения(Результат, Адрес, 
        ВыбранноеИмяФайла, ДополнительныеПараметры) экспорт
	//Сообщить("Процедура ОбработчикОкончанияПомещения");
	Если НЕ Результат Тогда
		//Сообщить("Процедура ОбработчикОкончанияПомещения --");
        Сообщить("Не удалось открыть обработку печати чека");
		возврат;
	КонецЕсли; 
	имяОбр = МенеджерОборудованияВызовСервера.ПолучитьИмяВнешнейОбработки(Адрес);	
	ПарамДляОбработкиПечатиЧека = Новый Структура;
	//передаём данные на основе которых можно было бы сформировать чека
	ПарамДляОбработкиПечатиЧека.Вставить("докСсылкаПоКоторомуПечатаемЧек", ДополнительныеПараметры.докСсылкаПоКоторомуПечатаемЧек);
	ПарамДляОбработкиПечатиЧека.Вставить("Команда", ДополнительныеПараметры.Команда);
	////передаём форму документа, что спразу же во внешней обработке её модифицировать (записать номер чека)
	//ПарамДляОбработкиПечатиЧека.Вставить("ФормаДокум", ДополнительныеПараметры.ФормаДокум);
	//при закрытии формы документа, форма печати чека тоже закроется
	ФормаОжиданияОтвета = ПолучитьФорму("ВнешняяОбработка."+имяОбр+".Форма.Форма", ПарамДляОбработкиПечатиЧека, ДополнительныеПараметры.ФормаДокум);
	//после закрытия формы печати чека - вызывает процедуру обработки результата печати чека
	ФормаОжиданияОтвета.ОписаниеОповещенияОЗакрытии = ДополнительныеПараметры.ОповещениеПриЗавершенииПечатиЧека;
	//открываем форму, которая нам напечатает чек
	ФормаОжиданияОтвета.Открыть();	
	//Сообщить("Процедура ОбработчикОкончанияПомещения --");
КонецПроцедуры //ОбработчикОкончанияПомещения

#КонецОбласти   

------------------------

Функция ПодключитьУстройство(ОбъектДрайвера, Параметры, ПараметрыПодключения, ВыходныеПараметры) Экспорт
................
	Для Каждого Параметр Из Параметры Цикл
		Если Лев(Параметр.Ключ, 2) = "P_" Тогда
			ЗначениеПараметра = Параметр.Значение;
			ИмяПараметра = Сред(Параметр.Ключ, 3);
			ОбъектДрайвера.УстановитьПараметр(ИмяПараметра, ЗначениеПараметра) 
		КонецЕсли;
	КонецЦикла;
	
	//ВстариваниеВКонфуУТ11 2024
	Ответ = ИСТИНА;
	Если ТипОборудования <> "ФискальныйРегистратор" тогда
		Ответ = ОбъектДрайвера.Подключить(ПараметрыПодключения.ИДУстройства);
	конецЕсли;
	//ВстариваниеВКонфуУТ11 2024 ---
	
	Если НЕ Ответ Тогда
		Результат = Ложь;
		ВыходныеПараметры.Добавить(999);
		ВыходныеПараметры.Добавить("");
		ОбъектДрайвера.ПолучитьОшибку(ВыходныеПараметры[1])
	Иначе
................
КонецФункции

-------------------

Функция ОтключитьУстройство(ОбъектДрайвера, Параметры, ПараметрыПодключения, ВыходныеПараметры) Экспорт
	
	Результат = Истина;
	ВыходныеПараметры = Новый Массив();	
	//ВстариваниеВКонфуУТ11 2024	
	ТипОборудования = ПараметрыПодключения.ТипОборудования;	
	Если ТипОборудования <> "ФискальныйРегистратор" тогда
		ОбъектДрайвера.Отключить(ПараметрыПодключения.ИДУстройства);
	конецЕсли;
	//ВстариваниеВКонфуУТ11 2024 ---
	
	Возврат Результат;
	
КонецФункции

-------------

Функция ВыполнитьКоманду(Команда, ВходныеПараметры = Неопределено, ВыходныеПараметры = Неопределено,
                         ОбъектДрайвера, Параметры, ПараметрыПодключения) Экспорт
	
	Результат = Истина;
	
	ВыходныеПараметры = Новый Массив();
	
	// ПРОЦЕДУРЫ И ФУНКЦИИ ОБЩИЕ ДЛЯ ВСЕХ ТИПОВ ДРАЙВЕРОВ
	
	//ВстариваниеВКонфуУТ11 2024
	ТипОборудования = ПараметрыПодключения.ТипОборудования;		
	Если ТипОборудования <> "ФискальныйРегистратор" тогда
		Если Команда = "ПолучитьВерсиюДрайвера" ИЛИ Команда = "GetVersion" или
			 Команда = "ПолучитьОписаниеДрайвера" ИЛИ Команда = "GetDescription" Тогда		
		 иначе
			 Возврат Результат;
		конецЕсли;
	конецЕсли;
	//ВстариваниеВКонфуУТ11 2024 ---

