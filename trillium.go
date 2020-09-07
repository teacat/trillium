package trillium

import (
	"errors"
	"net"
	"sync"
	"time"
)

const (
	// bitLenTime 是時間戳的位元長度。
	bitLenTime = 39
	// bitLenSequence 是順序位元長度。
	bitLenSequence = 8
	// bitLenWorkerID 是工作編號長度。
	bitLenWorkerID = 63 - bitLenTime - bitLenSequence
)

// trilliumTimeUnit 是 Trillium 的時間單位，這裡即是 `nsec` 亦為 `10 msec`。
const trilliumTimeUnit = 1e7

var (
	// ErrNoIP 表示找不到 IP 位置能夠作為工作編號使用。
	ErrNoIP = errors.New("trillium: no private ip address")
	// ErrOverTimeLimit 表示時間已經超過可用範圍。
	ErrOverTimeLimit = errors.New("trillium: over the time limit")
)

// Config 是初始化 Trillium 的設置檔。
type Config struct {
	// Since 是時間開始的 Unix 時間戳，全體服務應設置相同的時間並且就不得再次更改。
	Since time.Time
	// WorkerID 是此 Trillium 的工作編號，最高為 99,999。
	WorkerID uint16
}

// Trillium 是能在分佈式系統配發唯一編號的結構體。
type Trillium struct {
	// mutex 是同步鎖，能夠確保產生編號的時候不會因為多執行緒而導致衝突或重複碰撞。
	mutex *sync.Mutex
	// since 是時間開始的 Unix 時間戳，全體服務應設置相同的時間並且就不得再次更改。
	since int64
	// elapsedTime 是最後一次產生的時間，以此來確保可用編號的週期。
	elapsedTime int64
	// sequence 是目前時間內所可用的編號（累計直至下一個週期）。
	sequence uint16
	// workerID 是建立此 Trillium 時所配發的隨機編號，為了去中心化相依性而避免編號產生重複。
	workerID uint16
}

// DefaultConfig 會回傳一個預設的設定檔。
func DefaultConfig() *Config {
	return &Config{
		Since:    time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
		WorkerID: lower16BitPrivateIP(),
	}
}

// New 能夠建立新的 Trillium 來在分佈式系統中配發唯一不重複編號，
// `since` 是唯一編號時間戳的起始時間，可以自訂但之後就不應更改，
// 傳入 `0` 將會以預設的 `900288000` 作為基準。
func New(c *Config) *Trillium {
	t := new(Trillium)
	t.mutex = new(sync.Mutex)
	t.sequence = uint16(1<<bitLenSequence - 1)
	t.workerID = c.WorkerID
	t.since = toTrilliumTime(c.Since)
	return t
}

// Generate 會回傳新的唯一編號，此函式每秒可以產生 25,600 個唯一編號，
// 如果該秒的額度耗盡，將會自動推遲到下一秒才回傳。
func (t *Trillium) Generate() (uint64, error) {
	const maskSequence = uint16(1<<bitLenSequence - 1)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	current := currentElapsedTime(t.since)
	if t.elapsedTime < current {
		t.elapsedTime = current
		t.sequence = 0
	} else {
		t.sequence = (t.sequence + 1) & maskSequence
		if t.sequence == 0 {
			t.elapsedTime++
			overtime := t.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	return t.toID()
}

// toTrilliumTime 會將時間轉為 Trillium 格式。
func toTrilliumTime(t time.Time) int64 {
	return t.UTC().UnixNano() / trilliumTimeUnit
}

// currentElapsedTime 會回傳目前經過的時間以確保這個時間範圍內還有可用編號。
func currentElapsedTime(since int64) int64 {
	return toTrilliumTime(time.Now()) - since
}

// sleepTime 會回傳一個計算後的休息時間以讓 Trillium 重設資料。
func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%trilliumTimeUnit)*time.Nanosecond
}

// toID 會將位元資料整合後回傳一個 Trillium 編號。
func (t *Trillium) toID() (uint64, error) {
	if t.elapsedTime >= 1<<bitLenTime {
		return 0, ErrOverTimeLimit
	}

	return uint64(t.elapsedTime)<<(bitLenSequence+bitLenWorkerID) |
		uint64(t.sequence)<<bitLenWorkerID |
		uint64(t.workerID), nil
}

// privateIPv4 會盡可能地取得私人 IPv4 地址以做為工作編號使用。
func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, ErrNoIP
}

// isPrivateIPv4 會驗證一個 IP 地址是否為私人 IPv4。
func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

// lower16BitPrivateIP 會從私人 IPv4 中取得最低的 16 位元地址作為工作編號。
func lower16BitPrivateIP() uint16 {
	ip, err := privateIPv4()
	if err != nil {
		return 0
	}
	return uint16(ip[2])<<8 + uint16(ip[3])
}
