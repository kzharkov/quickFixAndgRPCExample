package main

import (
	"flag"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"os"
	"os/signal"
	"path"
	"quickFix/adapter"
	"quickFix/app"
)

func main() {
	address := flag.String("a", "127.0.0.1", "Address gRPC")
	flag.Parse()

	/*
		Здесь мы просто задаём путь к файлу конфигурации,
		тоже самое, если бы просто написали cfgFileName := "config/quickFix.cfg" (или "config\quickFix.cfg" в Windows).
		Поэтому тут и используется функция path.Join() - чтобы устранить платформозависимость
	*/
	cfgFileName := path.Join("config", "quickFix.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	/*
		Просто открываем файл и создаём объект файла cfg
	*/
	cfg, err := os.Open(cfgFileName)
	if err != nil {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}

	/*
		Функция читает настройки из нашего открытого файла
	*/
	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		fmt.Println("Error reading cfg,", err)
		return
	}

	/*
		logFactory - объект для логгирования
	*/
	logFactory := quickfix.NewScreenLogFactory()

	/*
		Наш объект для gRPC
	*/
	gRPC, err := adapter.NewClientGRPC(*address)

	/*
		Создаём объект lmax для интерфейса adapter
	*/
	lmax := adapter.NewLmaxAdapter(gRPC)

	/*
		Структура app принимает только интерфейс Adapter,
		а так как в структуре lmax имплементированы все функции Adapter,
		то мы можем подавать lmax, как аргумент для создания app.
	*/
	application := app.NewApplication(lmax)

	/*
		Ассоциируем нашу объект application с lmax адаптером
	*/
	lmax.SetApplication(application)
	appSettings.GlobalSettings().BoolSetting("Password")

	acceptor, err := quickfix.NewInitiator(application, quickfix.NewMemoryStoreFactory(), appSettings, logFactory)
	if err != nil {
		fmt.Printf("Unable to create Acceptor: %s\n", err)
		return
	}

	err = acceptor.Start()
	if err != nil {
		fmt.Printf("Unable to start Acceptor: %s\n", err)
		return
	}

	/*
		Запускаем наш хандлер
	*/
	//go func() {
	//	err := gRPC.DoStreamCommand()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	go func() {
		<-interrupt
		acceptor.Stop()
		os.Exit(0)
	}()
	select {}
}
