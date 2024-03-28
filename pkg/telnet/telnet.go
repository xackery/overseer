package telnet

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/xackery/overseer/pkg/flog"
	"github.com/ziutek/telnet"
)

var (
	mu            sync.Mutex
	onlineCount   int
	avgLevel      int
	popularClass  string
	nextRefresh   time.Time
	isInitialized bool
)

const (
	linebreak = "\n\r> "
)

var classNames = map[int]string{
	1:  "Warrior",
	2:  "Cleric",
	3:  "Paladin",
	4:  "Ranger",
	5:  "Shadowknight",
	6:  "Druid",
	7:  "Monk",
	8:  "Bard",
	9:  "Rogue",
	10: "Shaman",
	11: "Necromancer",
	12: "Wizard",
	13: "Magician",
	14: "Enchanter",
	15: "Beastlord",
	16: "Berserker",
}

// OnlineCount returns the number of online users.
func OnlineCount() int {
	mu.Lock()
	defer mu.Unlock()
	return onlineCount
	// if !isInitialized {
	// 	isInitialized = true
	// 	nextRefresh = time.Now().Add(3 * time.Second)
	// }
	// if time.Now().Before(nextRefresh) {
	// 	return onlineCount
	// }
	// nextRefresh = time.Now().Add(1 * time.Minute)
	// refreshStats()
	// return onlineCount
}

func AvgLevel() int {
	mu.Lock()
	defer mu.Unlock()
	return avgLevel
	// if !isInitialized {
	// 	isInitialized = true
	// 	nextRefresh = time.Now().Add(3 * time.Second)
	// }
	// if time.Now().Before(nextRefresh) {
	// 	return avgLevel
	// }
	// nextRefresh = time.Now().Add(1 * time.Minute)
	// refreshStats()
	// return avgLevel
}

func PopularClass() string {
	mu.Lock()
	defer mu.Unlock()
	return popularClass
	// if !isInitialized {
	// 	isInitialized = true
	// 	nextRefresh = time.Now().Add(3 * time.Second)
	// }
	// if time.Now().Before(nextRefresh) {
	// 	return popularClass
	// }
	// nextRefresh = time.Now().Add(1 * time.Minute)
	// refreshStats()
	// return popularClass
}

func refreshStats() {
	flog.Printf("[telnet] refreshing stats\n")
	conn, err := connect()
	if err != nil {
		flog.Printf("[telnet] connect: %s\n", err)
		return
	}

	onlineCount = 0

	type apiRespStruct struct {
		Data []struct {
			AccountID            int    `json:"account_id"`
			AccountName          string `json:"account_name"`
			Admin                int    `json:"admin"`
			Anon                 int    `json:"anon"`
			CharacterID          int    `json:"character_id"`
			Class                int    `json:"class"`
			ClientVersion        int    `json:"client_version"`
			Gm                   int    `json:"gm"`
			GuildID              int64  `json:"guild_id"`
			GuildRank            int    `json:"guild_rank"`
			GuildTributeOptIn    bool   `json:"guild_tribute_opt_in"`
			ID                   int    `json:"id"`
			Instance             int    `json:"instance"`
			IP                   int    `json:"ip"`
			IsLocalClient        bool   `json:"is_local_client"`
			Level                int    `json:"level"`
			Lfg                  bool   `json:"lfg"`
			LfgComments          string `json:"lfg_comments"`
			LfgFromLevel         int    `json:"lfg_from_level"`
			LfgMatchFilter       bool   `json:"lfg_match_filter"`
			LfgToLevel           int    `json:"lfg_to_level"`
			LoginserverAccountID int    `json:"loginserver_account_id"`
			LoginserverID        int    `json:"loginserver_id"`
			LoginserverName      string `json:"loginserver_name"`
			Name                 string `json:"name"`
			Online               int    `json:"online"`
			Race                 int    `json:"race"`
			Server               any    `json:"server"`
			TellsOff             int    `json:"tells_off"`
			WorldAdmin           int    `json:"world_admin"`
			Zone                 int    `json:"zone"`
		} `json:"data"`
		ExecutionTime string `json:"execution_time"`
		Method        string `json:"method"`
	}
	apiResp := &apiRespStruct{}
	resp, err := command(conn, "api get_client_list")
	if err != nil {
		flog.Printf("[telnet] api get_client_list: %s\n", err)
		return
	}

	flog.Printf("resp: %s\n", resp)
	err = json.Unmarshal([]byte(resp), &apiResp)
	if err != nil {
		flog.Printf("[telnet] get_client_list unmarshal: %s\n", err)
		return
	}

	flog.Printf("[telnet] api get_client_list: %s\n", resp)

	avgLevel := 0
	classes := make(map[int]int)

	for _, client := range apiResp.Data {
		if client.Online == 0 {
			continue
		}
		avgLevel += client.Level
		classes[client.Class]++
	}

	popularClass = "None"
	popularClassCount := 0
	for class, count := range classes {
		if count == 0 {
			continue
		}
		if count > popularClassCount {
			popularClassCount = count
			popularClass = classNames[class]
		}
	}

	onlineCount = len(apiResp.Data)
}

func connect() (*telnet.Conn, error) {
	var err error
	conn, err := telnet.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("set read deadline: %w", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("set write deadline: %w", err)
	}
	index := 0

	index, err = conn.SkipUntilIndex("Username:", "Connection established from localhost, assuming admin")
	if err != nil {
		return nil, fmt.Errorf("unexpected initial handshake: %w", err)
	}
	if index == 0 {
		return nil, fmt.Errorf("local bypass auth failed")
	}

	err = sendln(conn, "echo off")
	if err != nil {
		return nil, fmt.Errorf("echo off: %w", err)
	}

	err = sendln(conn, "acceptmessages off")
	if err != nil {
		return nil, fmt.Errorf("acceptmessages off: %w", err)
	}

	return conn, nil
}

// sendln sends a line to the telnet server.
func sendln(conn *telnet.Conn, s string) error {
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'

	flog.Printf("[telnet] sendln: %s\n", s)

	_, err := conn.Write(buf)
	if err != nil {
		return fmt.Errorf("sendLn: %s: %w", s, err)
	}
	return nil
}

// command sends a command to the telnet server and returns the output
func command(conn *telnet.Conn, cmd string) (string, error) {
	var err error

	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return "", fmt.Errorf("set read deadline: %w", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return "", fmt.Errorf("set write deadline: %w", err)
	}

	err = sendln(conn, cmd)
	if err != nil {
		return "", fmt.Errorf("sendln: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			flog.Printf("[telnet] panic in command: %s\n", r)
		}
	}()

	var data []byte
	var output string
	for {
		data, err = conn.ReadUntil(linebreak)
		if err != nil {
			return "", fmt.Errorf("read until: %w", err)
		}

		flog.Printf("[telnet] data: %s\n", string(data))

		output += string(data)

		if strings.Contains(output, linebreak) {
			output = strings.Replace(output, linebreak, "", 1)
			return output, nil
		}
	}
}
