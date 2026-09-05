package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte
var jwtSecretEphemeral bool

func init() {
	var err error
	jwtSecret, jwtSecretEphemeral, err = resolveJWTSecret(os.Getenv("JWT_SECRET"))
	if err != nil {
		panic(err)
	}
}

func resolveJWTSecret(value string) ([]byte, bool, error) {
	if value != "" {
		if len(value) < minSecretBytes {
			return nil, false, fmt.Errorf("JWT_SECRET must be at least %d bytes", minSecretBytes)
		}
		return []byte(value), false, nil
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, false, fmt.Errorf("generate JWT secret: %w", err)
	}
	return secret, true, nil
}

func resolveSetupToken(userCount int, value string) (string, bool, error) {
	if userCount > 0 {
		return "", false, nil
	}
	value = strings.TrimSpace(value)
	if value == "" {
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return "", false, fmt.Errorf("generate SETUP_TOKEN: %w", err)
		}
		return hex.EncodeToString(secret), true, nil
	}
	if len(value) < minSecretBytes {
		return "", false, fmt.Errorf("SETUP_TOKEN must be at least %d bytes", minSecretBytes)
	}
	return value, false, nil
}

func setupTokenMatches(expected, provided string) bool {
	if expected == "" {
		return false
	}
	expectedHash := sha256.Sum256([]byte(expected))
	providedHash := sha256.Sum256([]byte(provided))
	return subtle.ConstantTimeCompare(expectedHash[:], providedHash[:]) == 1
}

func generateToken(userID int64, username string) (string, error) {
	header := base64url([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, err := json.Marshal(struct {
		Sub  int64  `json:"sub"`
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}{
		Sub:  userID,
		Name: username,
		Exp:  time.Now().Add(72 * time.Hour).Unix(),
	})
	if err != nil {
		return "", err
	}
	payloadEnc := base64url(payload)
	sig := hmacSHA256(header+"."+payloadEnc, jwtSecret)
	return header + "." + payloadEnc + "." + sig, nil
}

func validateToken(token string) (int64, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, "", fmt.Errorf("invalid token")
	}
	expectedSig := hmacSHA256(parts[0]+"."+parts[1], jwtSecret)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return 0, "", fmt.Errorf("invalid signature")
	}
	payload, err := base64urlDecode(parts[1])
	if err != nil {
		return 0, "", err
	}
	var claims struct {
		Sub  int64  `json:"sub"`
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, "", err
	}
	if time.Now().Unix() > claims.Exp {
		return 0, "", fmt.Errorf("token expired")
	}
	return claims.Sub, claims.Name, nil
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// Crypto helpers
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

func hmacSHA256(data string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64url(h.Sum(nil))
}

func base64url(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func base64urlDecode(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func (d *DB) UserCount() int {
	var n int
	d.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&n)
	return n
}

func (d *DB) CreateUser(username, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	res, err := d.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

var errAdminAlreadyExists = errors.New("admin user already exists")
var errInvalidCredentials = errors.New("invalid username or password")

func (d *DB) CreateInitialUser(username, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	res, err := d.db.Exec(`
		INSERT INTO users (username, password_hash)
		SELECT ?, ?
		WHERE NOT EXISTS (SELECT 1 FROM users)
	`, username, string(hash))
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows != 1 {
		return 0, errAdminAlreadyExists
	}
	return res.LastInsertId()
}

var invalidUserPasswordHash = func() []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte("streamdock-invalid-user"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hash
}()

func (d *DB) VerifyUser(username, password string) (int64, error) {
	var id int64
	var hash string
	err := d.db.QueryRow("SELECT id, password_hash FROM users WHERE username=?", username).Scan(&id, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(invalidUserPasswordHash, []byte(password))
		return 0, errInvalidCredentials
	}
	if err != nil {
		return 0, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, errInvalidCredentials
	}
	return id, nil
}

type loginAttempt struct {
	failures     int
	firstFailure time.Time
	blockedUntil time.Time
	lastSeen     time.Time
}

type loginRateLimiter struct {
	mu         sync.Mutex
	attempts   map[string]loginAttempt
	maxEntries int
}

func newLoginRateLimiter() *loginRateLimiter {
	return newLoginRateLimiterWithLimit(maxTrackedLoginClients)
}

func newLoginRateLimiterWithLimit(maxEntries int) *loginRateLimiter {
	if maxEntries < 1 {
		maxEntries = 1
	}
	return &loginRateLimiter{
		attempts:   make(map[string]loginAttempt),
		maxEntries: maxEntries,
	}
}

func (l *loginRateLimiter) pruneExpired(now time.Time) {
	for client, attempt := range l.attempts {
		if now.Before(attempt.blockedUntil) {
			continue
		}
		if attempt.firstFailure.IsZero() || !now.Before(attempt.firstFailure.Add(loginFailureWindow)) {
			delete(l.attempts, client)
		}
	}
}

func (l *loginRateLimiter) evictLeastRecentlySeen() {
	var oldestClient string
	var oldestSeen time.Time
	for client, attempt := range l.attempts {
		seen := attempt.lastSeen
		if seen.IsZero() {
			seen = attempt.firstFailure
		}
		if oldestClient == "" || seen.Before(oldestSeen) {
			oldestClient = client
			oldestSeen = seen
		}
	}
	if oldestClient != "" {
		delete(l.attempts, oldestClient)
	}
}

func (l *loginRateLimiter) allow(client string, now time.Time) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.pruneExpired(now)
	attempt, ok := l.attempts[client]
	if !ok {
		return true, 0
	}
	attempt.lastSeen = now
	if now.Before(attempt.blockedUntil) {
		l.attempts[client] = attempt
		return false, attempt.blockedUntil.Sub(now)
	}
	l.attempts[client] = attempt
	return true, 0
}

func (l *loginRateLimiter) recordFailure(client string, now time.Time) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.pruneExpired(now)
	attempt, exists := l.attempts[client]
	if !exists && len(l.attempts) >= l.maxEntries {
		l.evictLeastRecentlySeen()
	}
	if attempt.firstFailure.IsZero() || now.Sub(attempt.firstFailure) >= loginFailureWindow {
		attempt = loginAttempt{firstFailure: now}
	}
	attempt.failures++
	attempt.lastSeen = now
	if attempt.failures >= maxLoginFailures {
		attempt.blockedUntil = now.Add(loginLockoutDuration)
	}
	l.attempts[client] = attempt
	if now.Before(attempt.blockedUntil) {
		return true, attempt.blockedUntil.Sub(now)
	}
	return false, 0
}

func (l *loginRateLimiter) reset(client string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, client)
}

var trustedProxyNets []*net.IPNet

func parseTrustedProxyCIDRs(value string) ([]*net.IPNet, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	var nets []*net.IPNet
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if !strings.Contains(part, "/") {
			if strings.Contains(part, ":") {
				part += "/128"
			} else {
				part += "/32"
			}
		}
		_, network, err := net.ParseCIDR(part)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR %q: %w", part, err)
		}
		nets = append(nets, network)
	}
	return nets, nil
}

func peerIP(r *http.Request) net.IP {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.ParseIP(r.RemoteAddr)
	}
	return net.ParseIP(host)
}

func peerIsTrusted(ip net.IP) bool {
	if ip == nil || len(trustedProxyNets) == 0 {
		return false
	}
	for _, network := range trustedProxyNets {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func parseForwardedClientIP(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	first := strings.TrimSpace(strings.Split(value, ",")[0])
	ip := net.ParseIP(first)
	if ip == nil {
		return ""
	}
	return ip.String()
}

func requestClientKey(r *http.Request) string {
	peer := peerIP(r)
	peerKey := "unknown"
	if peer != nil {
		peerKey = peer.String()
	} else if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil && host != "" {
		peerKey = host
	} else if r.RemoteAddr != "" {
		peerKey = r.RemoteAddr
	}
	if !peerIsTrusted(peer) {
		return peerKey
	}
	if realIP := parseForwardedClientIP(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	if forwarded := parseForwardedClientIP(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return forwarded
	}
	return peerKey
}

func (a *App) limiter() *loginRateLimiter {
	if a.loginLimiter == nil {
		a.loginLimiter = newLoginRateLimiter()
	}
	return a.loginLimiter
}

func (a *App) authRateLimitErr(w http.ResponseWriter, msg string, retryAfter time.Duration) {
	seconds := int(retryAfter / time.Second)
	if retryAfter%time.Second != 0 {
		seconds++
	}
	if seconds < 1 {
		seconds = 1
	}
	w.Header().Set("Retry-After", strconv.Itoa(seconds))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":               msg,
		"retry_after_seconds": seconds,
	})
}

func (a *App) setSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(sessionDuration),
		MaxAge:   int(sessionDuration.Seconds()),
		Secure:   requestIsHTTPS(r),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (a *App) clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		Secure:   requestIsHTTPS(r),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func sessionTokenFromRequest(r *http.Request) (token string, viaCookie bool) {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")), false
	}
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		return c.Value, true
	}
	return "", false
}

func (a *App) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, viaCookie := sessionTokenFromRequest(r)
		if token == "" {
			a.jsonErr(w, 401, "missing session")
			return
		}
		if _, _, err := validateToken(token); err != nil {
			a.jsonErr(w, 401, "token expired or invalid")
			return
		}
		if viaCookie && stateChangingMethod(r.Method) && !requestHasSameOrigin(r) {
			a.jsonErr(w, 403, "same-origin request required")
			return
		}
		next(w, r)
	}
}

// POST /api/auth/setup
func (a *App) handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		a.jsonErr(w, 405, "method not allowed")
		return
	}
	client := requestClientKey(r)
	if allowed, retryAfter := a.limiter().allow(client, time.Now()); !allowed {
		a.authRateLimitErr(w, "too many setup attempts; try again later", retryAfter)
		return
	}
	a.setupTokenMu.Lock()
	defer a.setupTokenMu.Unlock()
	if a.db.UserCount() > 0 {
		a.jsonErr(w, 400, "admin user already exists")
		return
	}
	var req struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		SetupToken string `json:"setup_token"`
	}
	if err := decodeJSONBody(w, r, &req); err != nil {
		a.jsonErr(w, 400, "invalid request")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || len(req.Username) > maxAdminUsernameChars || len(req.Password) < minAdminPasswordBytes || len(req.Password) > maxAdminPasswordBytes {
		a.jsonErr(w, 400, "username must be 1-64 characters and password must be 12-72 bytes")
		return
	}
	if a.setupToken == "" || !setupTokenMatches(a.setupToken, req.SetupToken) {
		if blocked, retryAfter := a.limiter().recordFailure(client, time.Now()); blocked {
			a.authRateLimitErr(w, "too many setup attempts; try again later", retryAfter)
			return
		}
		a.jsonErr(w, http.StatusForbidden, "invalid setup token")
		return
	}
	id, err := a.db.CreateInitialUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, errAdminAlreadyExists) {
			a.jsonErr(w, http.StatusConflict, errAdminAlreadyExists.Error())
			return
		}
		a.jsonErr(w, 500, "unable to create admin user")
		return
	}
	a.limiter().reset(client)
	token, err := generateToken(id, req.Username)
	if err != nil {
		a.jsonErr(w, 500, err.Error())
		return
	}
	a.setupToken = ""
	w.Header().Set("Cache-Control", "no-store")
	a.setSessionCookie(w, r, token)
	a.jsonOK(w, map[string]interface{}{"username": req.Username})
}

// POST /api/auth/login
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		a.jsonErr(w, 405, "method not allowed")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	client := requestClientKey(r)
	if allowed, retryAfter := a.limiter().allow(client, time.Now()); !allowed {
		a.authRateLimitErr(w, "too many login attempts; try again later", retryAfter)
		return
	}
	if err := decodeJSONBody(w, r, &req); err != nil {
		if blocked, retryAfter := a.limiter().recordFailure(client, time.Now()); blocked {
			a.authRateLimitErr(w, "too many login attempts; try again later", retryAfter)
			return
		}
		a.jsonErr(w, 400, "invalid request")
		return
	}
	username := strings.TrimSpace(req.Username)
	if username == "" || len(username) > maxAdminUsernameChars || req.Password == "" || len(req.Password) > maxAdminPasswordBytes {
		if blocked, retryAfter := a.limiter().recordFailure(client, time.Now()); blocked {
			a.authRateLimitErr(w, "too many login attempts; try again later", retryAfter)
			return
		}
		a.jsonErr(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	id, err := a.db.VerifyUser(username, req.Password)
	if err != nil {
		blocked, retryAfter := a.limiter().recordFailure(client, time.Now())
		if blocked {
			a.authRateLimitErr(w, "too many login attempts; try again later", retryAfter)
			return
		}
		a.jsonErr(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	a.limiter().reset(client)
	token, err := generateToken(id, username)
	if err != nil {
		a.jsonErr(w, 500, err.Error())
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	a.setSessionCookie(w, r, token)
	a.jsonOK(w, map[string]interface{}{"username": username})
}

// POST /api/auth/logout
func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		a.jsonErr(w, 405, "method not allowed")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	a.clearSessionCookie(w, r)
	a.jsonOK(w, map[string]bool{"logged_out": true})
}

// GET /api/auth/check
func (a *App) handleAuthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	needsSetup := a.db.UserCount() == 0
	authenticated := false
	username := ""
	if !needsSetup {
		if token, _ := sessionTokenFromRequest(r); token != "" {
			if _, sessionUsername, err := validateToken(token); err == nil {
				authenticated = true
				username = sessionUsername
			}
		}
	}
	a.jsonOK(w, map[string]interface{}{
		"needs_setup":          needsSetup,
		"mode":                 "single_admin",
		"jwt_secret_ephemeral": jwtSecretEphemeral,
		"setup_token_required": needsSetup,
		"authenticated":        authenticated,
		"username":             username,
	})
}
