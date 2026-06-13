// ZvA GAME (web edition) - Tugas 12 Cloud Computing
// Server Go yang menjalankan logika battle asli (Zhongli vs Azhdaha)
// dan menyajikannya sebagai web app yang siap deploy di PaaS (Render).

package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
)

//go:embed index.html
var indexHTML []byte

// State permainan dikirim bolak-balik dari browser (stateless di server).
type State struct {
	Zhp    int `json:"zhp"`
	Ahp    int `json:"ahp"`
	Energy int `json:"energy"`
}

type TurnRequest struct {
	State  State  `json:"state"`
	Action string `json:"action"` // "1", "2", atau "3"
}

type TurnResponse struct {
	State  State    `json:"state"`
	Log    []string `json:"log"`
	Status string   `json:"status"` // "ongoing", "win", "lose"
}

// processTurn menjalankan SATU giliran: aksi Zhongli, passive energy, lalu balasan Azhdaha.
// Ini port langsung dari loop di game terminal kamu.
func processTurn(s State, skil string) ([]string, State) {
	var logs []string
	add := func(line string) { logs = append(logs, line) }

	// Random sekali per giliran, dipakai bareng (sama seperti versi asli)
	randomInt := rand.Intn(500-200) + 200    // [200,500)
	randomInt2 := rand.Intn(1000-200) + 200   // [200,1000)
	randomHeal := rand.Intn(2000-50) + 50     // [50,2000)
	randomInt3 := rand.Intn(1000-500) + 500   // [500,1000)
	randomintz := rand.Intn(5)                // [0,5)
	randominta := rand.Intn(6)                // [0,6)

	var ULT int

	// ===== AKSI ZHONGLI =====
	switch skil {
	case "1":
		add("SPEAR!")
		spear := randomInt2 + 100
		if spear > 850 {
			add("✦ DMG: " + itoa(spear))
			spear += randomInt3 * 2
			add("✦ CRIT BONUS +" + itoa(randomInt3*2))
		}
		s.Ahp -= spear
		add("AZHDAHA: -" + itoa(spear) + " hp")
	case "2":
		add("DOMINUS SHIELD!")
		shield := randomInt
		s.Zhp += randomHeal + (shield * 2)
		add("ZHONGLI: Heal +" + itoa(randomHeal) + "  SHIELD ++" + itoa(shield))
	case "3":
		add("I WILL HAVE AN ORDER!!")
		ULT = randomintz
	}

	if skil == "3" {
		switch ULT {
		case 1:
			add("ORDER!! (Heal + Shield)")
			ult1 := 600 + (randomHeal * 2)
			add("✦ ZHONGLI EXTRA SHIELD: " + itoa(ult1))
			s.Zhp += ult1 + (randomHeal * 2)
			add("✦ ZHONGLI EXTRA HP: +" + itoa(randomHeal*2))
		case 2:
			add("PLANET BEFALL!!")
			ult2 := randomInt2 * 2
			if ult2 > 650 {
				add("✦ DMG: " + itoa(ult2))
				ult2 += randomInt3
				add("✦ CRIT BONUS +" + itoa(randomInt3))
			}
			s.Ahp -= ult2
			add("AZHDAHA: -" + itoa(ult2))
		case 0:
			add("DOMINANCE OF EARTH!!")
			ult3 := randomInt2 * 2
			if ult3 > 1200 {
				add("✦ DMG: " + itoa(ult3))
				add("✦ CRIT DMG: " + itoa(ult3*4))
				ult3 = ult3 * 5
			}
			s.Ahp -= ult3
			add("AZHDAHA: -" + itoa(ult3))
		case 4:
			add("✦ ULT CANCELLED")
			s.Energy++
		case 3:
			add("✦ ULT MISS")
			s.Energy++
		}
	}

	// ===== PASSIVE ENERGY (2 ULT gagal = pulih) =====
	if s.Energy == 2 {
		add("⚡ PASSIVE HEAL: 1800")
		add("⚡ PASSIVE ATK -> AZHDAHA: -900")
		s.Ahp -= 900
		s.Zhp += 1800
		s.Energy -= 2
	}

	// Kalau Azhdaha sudah kalah, dia tidak sempat membalas
	if s.Ahp <= 0 {
		return logs, s
	}

	// ===== BALASAN AZHDAHA =====
	azskil := randominta
	switch azskil {
	case 0:
		add("AZHDAHA: GEO")
		ad1 := randomInt3
		add("✦ GEO DMG: " + itoa(ad1))
		if ad1 > 800 {
			s.Ahp += randomHeal * 2
			add("✦ GEO HEAL BONUS: " + itoa(randomHeal*2))
			s.Ahp += randomInt3
			add("✦ GEO SHIELD BONUS: " + itoa(randomInt3))
		}
		s.Zhp -= ad1
		add("ZHONGLI: -" + itoa(ad1))
	case 1:
		add("AZHDAHA: CRYO")
		ad2 := randomInt2
		add("✦ CRYO DMG: " + itoa(ad2))
		if ad2 > 600 {
			ad2 += randomInt + 500
			add("✦ CRYO CRIT+: " + itoa(500+randomInt))
			s.Ahp += 500
			add("✦ AZHDAHA CRYO SHIELD: +500")
		}
		s.Zhp -= ad2
		add("ZHONGLI: -" + itoa(ad2))
	case 2:
		add("AZHDAHA: PYRO")
		ad3 := randomInt * 2
		add("✦ PYRO DMG: " + itoa(ad3))
		if ad3 > 650 {
			ad3 += randomInt3 * 2
			add("✦ PYRO BURN: " + itoa(randomInt3*2))
		}
		s.Zhp -= ad3
		add("ZHONGLI: -" + itoa(ad3))
	case 3:
		add("AZHDAHA: HYDRO")
		ad4 := randomInt2 + randomHeal
		add("✦ HYDRO DMG: " + itoa(ad4))
		if ad4 > 750 {
			ad4 += randomInt + randomHeal
			add("✦ HYDRO BONUS DMG: " + itoa(randomInt+randomHeal))
		}
		s.Zhp -= ad4
		add("ZHONGLI: -" + itoa(ad4))
	case 4:
		add("✦ AZHDAHA MISSED THE ATK")
		add("✦ ZHONGLI SHIELD COUNTER ATK!!")
		s.Ahp -= 1000
		add("AZHDAHA: -1000")
	case 5:
		add("AZHDAHA: ELECTRO")
		ad5 := randomInt3
		if ad5 >= 550 {
			add("✦ DMG: " + itoa(ad5))
			add("✦ CRIT DAMAGE: " + itoa(ad5+ad5))
			ad5 = ad5 * 3
			s.Zhp -= ad5
			add("ZHONGLI: -" + itoa(ad5))
		} else {
			ad5 = randomInt3 * 2
			s.Ahp += ad5
			add("AZHDAHA ++Heal: +" + itoa(ad5))
		}
	}

	return logs, s
}

func turnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req TurnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	logs, st := processTurn(req.State, req.Action)

	status := "ongoing"
	if st.Ahp <= 0 {
		status = "win"
	} else if st.Zhp <= 0 {
		status = "lose"
	}

	resp := TurnResponse{State: st, Log: logs, Status: status}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/turn", turnHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	// PaaS (Render) yang menentukan PORT lewat environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("ZvA server jalan di port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// itoa biar nggak perlu import strconv di banyak tempat
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
