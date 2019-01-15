package trillium

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	// maxSequence 是最大的流水編號。
	maxSequence int32 = 99999
	// maxRandomID 是最大的隨機服務編號。
	maxRandomID int = 99999
)

var (
	// from 是唯一編號時間戳的起始時間，可以自訂但之後就不應更改。
	from time.Time
)

// New 能夠建立新的 Trillium 來在分佈式系統中配發唯一不重複編號，
// `since` 是唯一編號時間戳的起始時間，可以自訂但之後就不應更改，
// 傳入 `0` 將會以預設的 `900288000` 作為基準。
func New(since int64) *Trillium {
	if since == 0 {
		since = 900288000
	}
	from = time.Unix(since, 0)

	// 替這個 Trillium 建立隨機編號避免與其他正在執行的服務發生編號碰撞。
	randomID := rand.Intn(maxRandomID)
	return &Trillium{
		randomID: int16(randomID),
		sequence: 0,
	}
}

// Trillium 是能在分佈式系統配發唯一編號的結構體。
type Trillium struct {
	// mutex 是同步鎖，能夠確保產生編號的時候不會因為多執行緒而導致衝突或重複碰撞。
	mutex sync.Mutex
	// timestamp 是最後一次產生的時間，以此來確保可用編號的週期。
	timestamp int64
	// sequence 是目前時間內所可用的編號（累計直至下一個週期）。
	sequence int32
	// randomID 是建立此 Trillium 時所配發的隨機編號，為了去中心化相依性而避免編號產生重複。
	randomID int16
}

// Generate 會回傳新的唯一編號，此函式每秒可以產生 100,000 個唯一編號，
// 如果該秒的額度耗盡，將會自動推遲到下一秒才回傳。
func (t *Trillium) Generate() *ID {
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
	// 將上次產生的時間設置為現在當下。
	t.timestamp = now
	t.mutex.Unlock()
	return &ID{
		Timestamp: time.Since(from).Seconds(),
		RandomID:  t.randomID,
		Sequence:  t.sequence,
	}
}

// ID 呈現了一個唯一編號的結構。
type ID struct {
	// Timestamp 是唯一編號建立時的時間戳。
	Timestamp float64
	// RandomID 是建立此唯一編號的服務隨機編號。
	RandomID int16
	// Sequence 是唯一編號的流水號。
	Sequence int32
}

// String 會回傳基於 `string` 型態的唯一編號。
func (i *ID) String() string {
	timestamp := int(time.Since(from).Seconds())
	randomID := fmt.Sprintf("%05d", i.RandomID)
	sequence := fmt.Sprintf("%05d", i.Sequence)

	return fmt.Sprintf("%d%s%s", timestamp, randomID, sequence)
}

// Int 會回傳基於 `int` 型態的唯一編號。
func (i *ID) Int() int {
	timestamp := int(time.Since(from).Seconds())
	randomID := fmt.Sprintf("%05d", i.RandomID)
	sequence := fmt.Sprintf("%05d", i.Sequence)

	id, _ := strconv.Atoi(fmt.Sprintf("%d%s%s", timestamp, randomID, sequence))
	return id
}
