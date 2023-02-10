package db

import (
	"log"

	"github.com/dgraph-io/badger/v2"
	"github.com/zippoxer/bow"
)

type Node struct {
	IP     string `bow:"key"` // primary key
	NAME   string
	STATUS string
}

type Connection struct {
	Id   string `bow:"key"`
	Node string
}

func openDB() *bow.DB {
	db, err := bow.Open("test",
		bow.SetBadgerOptions(badger.DefaultOptions("test").WithTruncate(true)))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var DB = openDB()

func GetNodes(query string) (Node, error) {
	var node Node
	var connection Connection

	err := DB.Bucket("tcp").Get(query, &connection)
	if err != nil {
		println("Error: ", err.Error())
		return node, err
	}
	println(connection.Id, " ", connection.Node)
	nodeErr := DB.Bucket("nodes").Get(connection.Node, &node)
	if nodeErr != nil {
		println("Error: ", nodeErr.Error())
		return node, nodeErr
	}
	println(node.IP, " ", node.NAME, " ", node.STATUS, " ", node.IP)
	return node, nil
}

func SetNode(data *Node) (bool, error) {
	err := DB.Bucket("nodes").Put(data)
	if err != nil {
		println("Error: ", err.Error())
		return false, err
	}

	return true, nil

}

func SetConnection(data *Connection) error {

	err := DB.Bucket("tcp").Put(data)
	if err != nil {
		println("Error: ", err.Error())
		return err
	}
	return nil

}

func ErrorNode(data *Node) (bool, error) {

	err := DB.Bucket("nodes").Delete(data.IP)
	if err != nil {
		println("Error: ", err.Error())
		return false, err
	}
	iter := DB.Bucket("tcp").Iter()
	defer iter.Close()
	var connection Connection
	for iter.Next(&connection) {
		if connection.Node == data.IP {
			err := DB.Bucket("tcp").Delete(connection.Id)
			if err != nil {
				println("Error: ", err.Error())
				return false, err
			}
		}
	}
	if iter.Err() != nil {

		return true, iter.Err()

	}

	return true, nil

}
