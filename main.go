package main

import (
	"fmt"
	"time"

	pbe "github.com/withlin/canal-go/protocol/entry"

	"github.com/golang/protobuf/proto"
	"github.com/withlin/canal-go/client"
)

func main() {
	address := "127.0.0.1"
	port := 11111
	username := ""
	password := ""
	destination := "example"
	soTimeout := int32(60000)
	idleTimeout := int32(60 * 60 * 1000)

	connector := client.NewSimpleCanalConnector(address, port, username, password, destination, soTimeout, idleTimeout)
	err := connector.Connect()
	if err != nil {
		panic(err)
	}

	// subscribe with db name and table name, support regex
	//err = connector.Subscribe(".*\\..*")
	err = connector.Subscribe("testdb.test")
	if err != nil {
		fmt.Printf("connector.Subscribe failed, err:%v\n", err)
		panic(err)
	}

	for {
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			fmt.Printf("connector.Get failed, err:%v\n", err)
			continue
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(time.Second)
			fmt.Println("no events")
			continue
		}
		handleEntry(message.Entries)
	}
}

func handleEntry(entries []pbe.Entry) {
	for _, entry := range entries {
		// ignore transaction begin and end
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}

		rowChange := new(pbe.RowChange)
		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		if err != nil {
			fmt.Printf("proto.Unmarshal failed, err:%v\n", err)
		}
		if rowChange == nil {
			continue
		}

		header := entry.GetHeader()
		fmt.Printf("listen binlog change: [%s : %d], [db=%s, table=%s], eventType: %s\n",
			header.GetLogfileName(),
			header.GetLogfileOffset(),
			header.GetSchemaName(),
			header.GetTableName(),
			header.GetEventType(),
		)

		// check if is DDL
		if rowChange.GetIsDdl() {
			fmt.Printf("DDL: %v\n", rowChange.GetSql())
		}

		// event type：insert/update/delete
		eventType := rowChange.GetEventType()
		for _, rowData := range rowChange.GetRowDatas() {
			if eventType == pbe.EventType_DELETE {
				handleColumn(rowData.GetBeforeColumns())
			} else if eventType == pbe.EventType_INSERT || eventType == pbe.EventType_UPDATE {
				handleColumn(rowData.GetAfterColumns())
			} else {
				fmt.Println("---before---")
				handleColumn(rowData.GetBeforeColumns())
				fmt.Println("---after---")
				handleColumn(rowData.GetAfterColumns())
			}
		}
	}
}

func handleColumn(columns []*pbe.Column) {
	for _, col := range columns {
		if col.GetUpdated() {
			fmt.Printf("%s: %s  update=%v\n", col.GetName(), col.GetValue(), col.GetUpdated())
		}
	}
}
