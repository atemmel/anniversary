package common

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"sync"
)

const DataLength = 4096

type Client struct {
	conf ClientConfig
	rw *bufio.ReadWriter
	conn net.Conn
	playerMap PlayerMap

	Active bool
	spinData *SpinMessage
}

type PlayerMap struct {
	players map[int]Player
	mutex sync.Mutex
}

const (
	PlayerMsg = 0
	SpinMsg = 1
)

type ClientMessage struct {
	MsgId int
	Player *Player
	SpinMsg *SpinMessage
}

type SpinMessage struct {
	Strings []string
	Offset int
}

type JoinMessage struct {
	Name string
	Id int
}

func CreateClient() Client {
	return Client{
		ClientConfig{},
		nil,
		nil,
		PlayerMap{
			make(map[int]Player),
			sync.Mutex{},
		},
		false,
		nil,
	}
}

func (c *Client) Connect() *JoinMessage {
	log.Println("Attempting to connect to server...")
	var err error

	c.conf, err = ReadClientConfig()
	if err != nil {
		log.Println("Could not read client config")
		return nil
	}

	c.conn, err = net.Dial("tcp", c.conf.ServerUrl + ":" + c.conf.ServerPort)
	if err != nil {
		log.Println("Connection failed")
		log.Println(err)
		return nil
	}

	log.Println("Connection succeeded!")
	c.rw = bufio.NewReadWriter(
		bufio.NewReader(c.conn),
		bufio.NewWriter(c.conn),
	)

	jmsg := &JoinMessage{}
	jmsg.Name = c.conf.DiscordName
	b, _ := json.Marshal(jmsg)
	b = append(b, '\n')
	_, err = c.rw.Write(b)
	if err != nil {
		log.Println(err)
		return nil
	}
	c.rw.Flush()

	data, err := c.rw.ReadBytes('\n')
	if err != nil {
		log.Println("Could not be given an id")
		return nil
	}

	data = data[:len(data) - 1]	// Discard newline byte
	//id, err := strconv.Atoi(string(data))
	err = json.Unmarshal(data, jmsg)
	if err != nil {
		log.Println("Id given (", string(data), ") was not valid")
		return nil
	}
	c.Active = true
	return jmsg
}

func (c *Client) WritePlayer(player *Player) {
	pmsg := ClientMessage{}
	pmsg.MsgId = PlayerMsg
	pmsg.Player = player
	data, _ := json.Marshal(&pmsg)
	data = append(data, '\n')
	c.rw.Write(data)
	c.rw.Flush()
}

func (c *Client) ReadPlayer() {
	for {
		bytes, err := c.rw.ReadBytes('\n')
		if err != nil {
			//TODO Identify which errors should be ignored and which 
			//     errors should abort the connection
			log.Println("Could not read message:", err)
			c.Active = false
			return
		} else {
			msg := &ClientMessage{}
			err := json.Unmarshal(bytes, msg)
			if err == nil {
				if msg.MsgId == PlayerMsg && msg.Player != nil {
					player := msg.Player
					c.updatePlayer(player)
				} else if msg.MsgId == SpinMsg && msg.SpinMsg != nil {
					//log.Println("Spinmsg rec", msg.SpinMsg)
					c.spinData = msg.SpinMsg
				} else {
					log.Println("Foreign message")
				}
			} else {
				log.Println(err)
				log.Println("Recieved:", string(bytes))
			}
		}
	}
}

func (c *Client) updatePlayer(player *Player) {
	c.playerMap.mutex.Lock()
	if !player.Connected {	// He disconnected
		delete(c.playerMap.players, player.Id)
		log.Println("Player", player.Id, "disconnected")
	} else {
		c.playerMap.players[player.Id] = *player
	}
	c.playerMap.mutex.Unlock()
}

func (c *Client) Disconnect() {
	if !c.Active {
		return
	}
	log.Println("Disconnecting...")
	c.Active = false
	c.rw.WriteByte(0)
	c.conn.Close()
}
