package trillium

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	maxSequence int32 = 99999
)

var (
	from time.Time
)

// New 能夠建立新的 Trillium 來在分佈式系統中配發唯一不重複編號。
func New(since int64) *Trillium {
	if since == 0 {
		since = 900288000
	}
	from = time.Unix(since, 0)

	// 替這個 Trillium 建立隨機編號避免與其他正在執行的服務發生編號碰撞。
	randomID := rand.Intn(99999)
	return &Trillium{
		randomID: int16(randomID),
		sequence: 0,
	}
}

// Trillium 是能在分佈式系統配發唯一編號的結構體。
type Trillium struct {
	// mutex 是同步鎖，能夠確保產生編號的時候不會因為多執行緒而導致衝突或重複碰撞。
	mutex sync.Mutex
	//
	timestamp int64
	// sequence 是目前時間內所可用的編號（累計直至下一個週期）。
	sequence int32
	// randomID 是建立此 Trillium 時所配發的隨機編號，為了去中心化相依性而避免編號產生重複。
	randomID int16
}

func (t *Trillium) Generate() int {
	t.mutex.Lock()
	// 取得產生唯一編號當下的時間。
	now := time.Now().Unix()
	// 如果上次產生的時間與現在完全一模一樣。
	if t.timestamp == now {
		// 而且編號額度已經達到上限的話。
		if t.sequence == maxSequence {
			// 就強迫等待到下一秒才繼續配發唯一編號。
			for t.timestamp == now {
				now = time.Now().Unix()
			}
			// 重新設置編號的配發額度。
			t.sequence = 0
		} else {
			// 否則就直接將配發額度遞增。
			t.sequence++
		}
	}

	// 組成 Trillium 的整體編號，並且在長度不足的時候自動補零。
	timestamp := int(time.Since(from).Seconds())
	randomID := fmt.Sprintf("%05d", t.randomID)
	sequence := fmt.Sprintf("%05d", t.sequence)
	id, _ := strconv.Atoi(fmt.Sprintf("%d%s%s", timestamp, randomID, sequence))

	// 將上次產生的時間設置為現在當下。
	t.timestamp = now

	t.mutex.Unlock()
	return id
}
