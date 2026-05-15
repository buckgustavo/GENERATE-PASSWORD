package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
)

const (
	UPPER = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LOWER = "abcdefghijklmnopqrstuvwxyz"
	NUMS  = "0123456789"
	SYMS  = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

type GenerateRequest struct {
	Length    int  `json:"length"`
	Uppercase bool `json:"uppercase"`
	Lowercase bool `json:"lowercase"`
	Numbers   bool `json:"numbers"`
	Symbols   bool `json:"symbols"`
}

type GenerateResponse struct {
	Password string `json:"password"`
	Entropy  int    `json:"entropy"`
	PoolSize int    `json:"poolSize"`
	Error    string `json:"error,omitempty"`
}

func secureRandom(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func generatePassword(req GenerateRequest) (string, int, error) {
	charset := ""
	var sets []string

	if req.Uppercase {
		charset += UPPER
		sets = append(sets, UPPER)
	}
	if req.Lowercase {
		charset += LOWER
		sets = append(sets, LOWER)
	}
	if req.Numbers {
		charset += NUMS
		sets = append(sets, NUMS)
	}
	if req.Symbols {
		charset += SYMS
		sets = append(sets, SYMS)
	}

	if charset == "" {
		return "", 0, fmt.Errorf("select at least one character set")
	}

	if req.Length < 4 || req.Length > 64 {
		return "", 0, fmt.Errorf("length must be between 4 and 64")
	}

	pw := make([]byte, 0, req.Length)

	for _, s := range sets {
		idx, err := secureRandom(len(s))
		if err != nil {
			return "", 0, err
		}
		pw = append(pw, s[idx])
	}

	for len(pw) < req.Length {
		idx, err := secureRandom(len(charset))
		if err != nil {
			return "", 0, err
		}
		pw = append(pw, charset[idx])
	}

	for i := len(pw) - 1; i > 0; i-- {
		j, err := secureRandom(i + 1)
		if err != nil {
			return "", 0, err
		}
		pw[i], pw[j] = pw[j], pw[i]
	}

	poolSize := len(charset)
	entropy := int(float64(req.Length) * math.Log2(float64(poolSize)))

	return string(pw), entropy, nil
}

func logBase2(x float64) float64 {
	result := 0.0
	for x >= 2.0 {
		x /= 2.0
		result++
	}
	return result + (x - 1.0)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(GenerateResponse{Error: "invalid request"})
		return
	}

	pw, entropy, err := generatePassword(req)
	if err != nil {
		json.NewEncoder(w).Encode(GenerateResponse{Error: err.Error()})
		return
	}

	poolSize := 0
	if req.Uppercase {
		poolSize += len(UPPER)
	}
	if req.Lowercase {
		poolSize += len(LOWER)
	}
	if req.Numbers {
		poolSize += len(NUMS)
	}
	if req.Symbols {
		poolSize += len(SYMS)
	}

	json.NewEncoder(w).Encode(GenerateResponse{
		Password: pw,
		Entropy:  entropy,
		PoolSize: poolSize,
	})
}

const htmlPage = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>PassForge — Go Edition</title>
<link href="https://fonts.googleapis.com/css2?family=Bebas+Neue&family=IBM+Plex+Mono:wght@400;500;600&display=swap" rel="stylesheet"/>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0b0c0f;--surface:#13151a;--border:#2a2d35;--accent:#e8ff47;
  --accent2:#ff6b35;--text:#e2e4ec;--muted:#5a5f72;
  --weak:#ff4444;--fair:#ff8c00;--good:#ffcc00;--strong:#7cff6b;--vstrong:#e8ff47;
}
body{background:var(--bg);color:var(--text);font-family:'IBM Plex Mono',monospace;
  min-height:100vh;display:flex;align-items:center;justify-content:center;
  padding:2rem 1rem;position:relative;overflow-x:hidden;}
body::before{content:'';position:fixed;inset:0;
  background:repeating-linear-gradient(0deg,transparent,transparent 40px,rgba(232,255,71,0.012) 40px,rgba(232,255,71,0.012) 41px),
  repeating-linear-gradient(90deg,transparent,transparent 40px,rgba(232,255,71,0.012) 40px,rgba(232,255,71,0.012) 41px);
  pointer-events:none;z-index:0;}
.glow-blob{position:fixed;width:500px;height:500px;border-radius:50%;
  background:radial-gradient(circle,rgba(232,255,71,0.04) 0%,transparent 70%);
  top:-100px;right:-100px;pointer-events:none;z-index:0;}
.container{position:relative;z-index:1;width:100%;max-width:520px;}
header{margin-bottom:2rem;}
.logo{font-family:'Bebas Neue',sans-serif;font-size:64px;letter-spacing:4px;
  color:var(--accent);line-height:1;position:relative;display:inline-block;}
.logo::after{content:'go edition';position:absolute;bottom:10px;right:-72px;
  font-family:'IBM Plex Mono',monospace;font-size:10px;color:var(--muted);letter-spacing:1px;}
.tagline{font-size:11px;color:var(--muted);letter-spacing:3px;text-transform:uppercase;margin-top:4px;}
.go-badge{display:inline-flex;align-items:center;gap:6px;margin-top:8px;
  background:rgba(0,173,216,0.1);border:1px solid rgba(0,173,216,0.3);
  border-radius:3px;padding:3px 8px;font-size:10px;color:#00add8;letter-spacing:2px;}
.pw-box{background:var(--surface);border:1px solid var(--border);border-radius:4px;
  padding:20px;margin-bottom:16px;position:relative;overflow:hidden;}
.pw-box::before{content:'OUTPUT';position:absolute;top:8px;left:16px;
  font-size:9px;letter-spacing:3px;color:var(--muted);}
.pw-display{font-size:22px;font-weight:600;letter-spacing:3px;color:var(--accent);
  word-break:break-all;line-height:1.5;margin-top:16px;min-height:36px;transition:opacity 0.1s;}
.pw-display.flash{animation:flash-pw 0.15s ease;}
@keyframes flash-pw{0%{opacity:0;transform:translateY(-4px)}100%{opacity:1;transform:translateY(0)}}
.strength-row{display:flex;align-items:center;gap:12px;margin-top:16px;}
.strength-bars{display:flex;gap:4px;flex:1;}
.strength-bar{height:4px;flex:1;border-radius:2px;background:var(--border);transition:background 0.3s ease;}
.strength-label{font-size:10px;letter-spacing:2px;text-transform:uppercase;min-width:70px;
  text-align:right;color:var(--muted);transition:color 0.3s;}
.copy-btn{width:100%;background:transparent;border:1px solid var(--accent);color:var(--accent);
  font-family:'IBM Plex Mono',monospace;font-size:12px;letter-spacing:3px;text-transform:uppercase;
  padding:12px;border-radius:4px;cursor:pointer;transition:background 0.15s,color 0.15s;margin-bottom:20px;}
.copy-btn:hover{background:rgba(232,255,71,0.08)}
.copy-btn:active{background:rgba(232,255,71,0.15);transform:scale(0.99)}
.copy-btn.copied{background:var(--accent);color:var(--bg);}
.panel{background:var(--surface);border:1px solid var(--border);border-radius:4px;padding:20px;margin-bottom:16px;}
.panel-label{font-size:9px;letter-spacing:3px;text-transform:uppercase;color:var(--muted);margin-bottom:16px;display:block;}
.length-row{display:flex;align-items:center;gap:16px;}
.length-val{font-family:'Bebas Neue',sans-serif;font-size:42px;color:var(--accent);min-width:52px;line-height:1;}
.slider-wrap{flex:1}
input[type=range]{-webkit-appearance:none;width:100%;height:4px;background:var(--border);border-radius:2px;outline:none;cursor:pointer;}
input[type=range]::-webkit-slider-thumb{-webkit-appearance:none;width:18px;height:18px;border-radius:50%;
  background:var(--accent);border:2px solid var(--bg);cursor:pointer;transition:transform 0.1s;}
input[type=range]::-webkit-slider-thumb:hover{transform:scale(1.2)}
input[type=range]::-moz-range-thumb{width:18px;height:18px;border-radius:50%;
  background:var(--accent);border:2px solid var(--bg);cursor:pointer;}
.slider-limits{display:flex;justify-content:space-between;font-size:10px;color:var(--muted);margin-top:6px;letter-spacing:1px;}
.options-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px;}
.opt{display:flex;align-items:center;gap:10px;cursor:pointer;padding:10px 12px;
  border:1px solid var(--border);border-radius:4px;transition:border-color 0.15s,background 0.15s;user-select:none;}
.opt:hover{border-color:var(--muted);background:rgba(255,255,255,0.02)}
.opt.checked{border-color:var(--accent);background:rgba(232,255,71,0.05)}
.opt input[type=checkbox]{display:none}
.opt-box{width:16px;height:16px;border:1px solid var(--muted);border-radius:2px;
  display:flex;align-items:center;justify-content:center;flex-shrink:0;transition:border-color 0.15s,background 0.15s;}
.opt.checked .opt-box{border-color:var(--accent);background:var(--accent);}
.opt-check{font-size:10px;color:var(--bg);opacity:0;transition:opacity 0.15s;}
.opt.checked .opt-check{opacity:1}
.opt-text{font-size:11px;letter-spacing:1px;color:var(--muted);transition:color 0.15s;line-height:1.3;}
.opt.checked .opt-text{color:var(--text)}
.opt-sample{font-size:9px;color:var(--muted);margin-top:2px;letter-spacing:1px;}
.gen-btn{width:100%;background:var(--accent);border:none;color:var(--bg);
  font-family:'Bebas Neue',sans-serif;font-size:24px;letter-spacing:4px;padding:16px;
  border-radius:4px;cursor:pointer;transition:transform 0.1s,box-shadow 0.15s;
  box-shadow:0 0 0 0 rgba(232,255,71,0);display:flex;align-items:center;justify-content:center;gap:10px;}
.gen-btn:hover{transform:translateY(-1px);box-shadow:0 4px 20px rgba(232,255,71,0.2);}
.gen-btn:active{transform:scale(0.99);box-shadow:none}
.gen-btn:disabled{opacity:0.5;cursor:not-allowed;transform:none;}
.gen-icon{font-size:20px}
.warning{font-size:10px;color:var(--accent2);letter-spacing:2px;text-align:center;margin-top:6px;min-height:14px;}
.entropy{font-size:10px;color:var(--muted);letter-spacing:2px;margin-top:8px;text-align:right;}
.entropy span{color:var(--text)}
</style>
</head>
<body>
<div class="glow-blob"></div>
<div class="container">
  <header>
    <div class="logo">PassForge</div>
    <div class="tagline">cryptographic password generator</div>
    <div class="go-badge">⚙ powered by Go crypto/rand Code by G.Buck</div>
  </header>

  <div class="pw-box">
    <div class="pw-display" id="pw-display">click generate</div>
    <div class="strength-row">
      <div class="strength-bars">
        <div class="strength-bar" id="sb1"></div>
        <div class="strength-bar" id="sb2"></div>
        <div class="strength-bar" id="sb3"></div>
        <div class="strength-bar" id="sb4"></div>
        <div class="strength-bar" id="sb5"></div>
      </div>
      <div class="strength-label" id="strength-label">---</div>
    </div>
    <div class="entropy" id="entropy-info"></div>
  </div>

  <button class="copy-btn" id="copy-btn" onclick="copyPw()">[ copy to clipboard ]</button>

  <div class="panel">
    <span class="panel-label">password length</span>
    <div class="length-row">
      <div class="length-val" id="len-val">16</div>
      <div class="slider-wrap">
        <input type="range" id="len-slider" min="4" max="64" value="16" oninput="onLenChange(this.value)"/>
        <div class="slider-limits"><span>4</span><span>64</span></div>
      </div>
    </div>
  </div>

  <div class="panel">
    <span class="panel-label">character sets</span>
    <div class="options-grid">
      <label class="opt checked" id="opt-upper">
        <input type="checkbox" checked onchange="toggleOpt(this,'opt-upper')"/>
        <div class="opt-box"><span class="opt-check">✓</span></div>
        <div><div class="opt-text">Uppercase</div><div class="opt-sample">A B C ... Z</div></div>
      </label>
      <label class="opt checked" id="opt-lower">
        <input type="checkbox" checked onchange="toggleOpt(this,'opt-lower')"/>
        <div class="opt-box"><span class="opt-check">✓</span></div>
        <div><div class="opt-text">Lowercase</div><div class="opt-sample">a b c ... z</div></div>
      </label>
      <label class="opt checked" id="opt-nums">
        <input type="checkbox" checked onchange="toggleOpt(this,'opt-nums')"/>
        <div class="opt-box"><span class="opt-check">✓</span></div>
        <div><div class="opt-text">Numbers</div><div class="opt-sample">0 1 2 ... 9</div></div>
      </label>
      <label class="opt" id="opt-syms">
        <input type="checkbox" onchange="toggleOpt(this,'opt-syms')"/>
        <div class="opt-box"><span class="opt-check">✓</span></div>
        <div><div class="opt-text">Symbols</div><div class="opt-sample">! @ # $ % &</div></div>
      </label>
    </div>
  </div>

  <div class="warning" id="warning"></div>

  <button class="gen-btn" id="gen-btn" onclick="generate()">
    GENERATE PASSWORD
  </button>
</div>

<script>
var currentPw = '';

function toggleOpt(cb, id) {
  var el = document.getElementById(id);
  if (cb.checked) el.classList.add('checked');
  else el.classList.remove('checked');
}

function onLenChange(v) {
  document.getElementById('len-val').textContent = v;
}

async function generate() {
  var btn = document.getElementById('gen-btn');
  btn.disabled = true;
  btn.innerHTML = '<span class="gen-icon">⏳</span> GENERATING...';

  document.getElementById('warning').textContent = '';

  var payload = {
    length:    parseInt(document.getElementById('len-slider').value),
    uppercase: document.querySelector('#opt-upper input').checked,
    lowercase: document.querySelector('#opt-lower input').checked,
    numbers:   document.querySelector('#opt-nums input').checked,
    symbols:   document.querySelector('#opt-syms input').checked,
  };

  try {
    var res = await fetch('/generate', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    });
    var data = await res.json();

    if (data.error) {
      document.getElementById('warning').textContent = '! ' + data.error;
      return;
    }

    currentPw = data.password;
    var disp = document.getElementById('pw-display');
    disp.textContent = data.password;
    disp.classList.remove('flash');
    void disp.offsetWidth;
    disp.classList.add('flash');

    updateStrength(data.entropy, data.poolSize, payload.length);

    var cb = document.getElementById('copy-btn');
    cb.classList.remove('copied');
    cb.textContent = '[ copy to clipboard ]';
  } catch(e) {
    document.getElementById('warning').textContent = '! server error — is the Go server running?';
  } finally {
    btn.disabled = false;
    btn.innerHTML =  'GENERATE PASSWORD';
  }
}

function updateStrength(bits, poolSize, len) {
  var bars = document.querySelectorAll('.strength-bar');
  var label = document.getElementById('strength-label');
  var info = document.getElementById('entropy-info');

  var level, color, text;
  if (bits < 28)       { level=1; color='var(--weak)';    text='WEAK'; }
  else if (bits < 50)  { level=2; color='var(--fair)';    text='FAIR'; }
  else if (bits < 72)  { level=3; color='var(--good)';    text='GOOD'; }
  else if (bits < 100) { level=4; color='var(--strong)';  text='STRONG'; }
  else                 { level=5; color='var(--vstrong)'; text='FORTRESS'; }

  bars.forEach(function(b, i) {
    b.style.background = i < level ? color : 'var(--border)';
  });
  label.textContent = text;
  label.style.color = color;
  info.innerHTML = 'entropy <span>' + bits + ' bits</span> &nbsp;|&nbsp; pool <span>' + poolSize + ' chars</span>';
}

function copyPw() {
  if (!currentPw) return;
  navigator.clipboard.writeText(currentPw).then(function() {
    var cb = document.getElementById('copy-btn');
    cb.classList.add('copied');
    cb.textContent = '[ ✓ copied! ]';
    setTimeout(function() {
      cb.classList.remove('copied');
      cb.textContent = '[ copy to clipboard ]';
    }, 2000);
  });
}

generate();
</script>
</body>
</html>`

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlPage)
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/generate", handleGenerate)

	port := "8080"
	fmt.Printf("PassForge (Go Edition) running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
