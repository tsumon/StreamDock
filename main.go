package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"streamdock/web"
)

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// UA Profiles 闂?only 3 modes
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// JWT helpers (simple HMAC-SHA256)
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

var appVersion = "v1.5.0"

const (
	sessionCookieName      = "streamdock_session"
	defaultDBFile          = "streamdock.db"
	legacyDBFile           = "meridian.db"
	backupFormatVersion    = "streamdock-v1"
	legacyBackupVersion    = "meridian-v1"
	backupDownloadName     = "streamdock_backup.json"
	sessionDuration        = 72 * time.Hour
	maxJSONBodyBytes       = 512 << 10
	maxLoginFailures       = 5
	maxTrackedLoginClients = 10000
	loginFailureWindow     = time.Minute
	loginLockoutDuration   = time.Minute
	minSecretBytes         = 32
	minAdminPasswordBytes  = 12
	maxAdminPasswordBytes  = 72
	maxAdminUsernameChars  = 64
	panelCSP               = "default-src 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'"
)

func envBool(name string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(name))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func panelListenAddr(bindAddress string, port int) (string, error) {
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("panel port must be between 1 and 65535, got %d", port)
	}
	bindAddress = strings.TrimSpace(bindAddress)
	if bindAddress == "" {
		bindAddress = "127.0.0.1"
	}
	if net.ParseIP(bindAddress) == nil {
		return "", fmt.Errorf("PANEL_BIND_ADDR must be an IP address, got %q", bindAddress)
	}
	return net.JoinHostPort(bindAddress, strconv.Itoa(port)), nil
}

func panelBindIsLoopback(bindAddress string) bool {
	bindAddress = strings.TrimSpace(bindAddress)
	if bindAddress == "" {
		return true
	}
	ip := net.ParseIP(bindAddress)
	return ip != nil && ip.IsLoopback()
}

// Minimal JWT 闂?no external dependency

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// Database
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?
type DB struct {
	db *sql.DB
}

func openDB(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	d := &DB{db: sqlDB}
	if err := d.migrate(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DB) Close() { d.db.Close() }

func (d *DB) migrate() error {
	_, err := d.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS sites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		listen_port INTEGER NOT NULL UNIQUE,
		target_url TEXT NOT NULL,
		playback_target_url TEXT NOT NULL DEFAULT '',
		stream_hosts TEXT NOT NULL DEFAULT '[]',
		ua_mode TEXT DEFAULT 'infuse',
		enabled INTEGER DEFAULT 1,
		traffic_quota BIGINT DEFAULT 0,
		traffic_used BIGINT DEFAULT 0,
		speed_limit INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS traffic_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site_id INTEGER NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
		bytes_in BIGINT DEFAULT 0,
		bytes_out BIGINT DEFAULT 0,
		recorded_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_traffic_site_time ON traffic_logs(site_id, recorded_at);
	`)
	if err != nil {
		return err
	}

	var hasPlaybackTargetColumn int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sites') WHERE name='playback_target_url'").Scan(&hasPlaybackTargetColumn); err != nil {
		return err
	}
	if hasPlaybackTargetColumn == 0 {
		if _, err := d.db.Exec("ALTER TABLE sites ADD COLUMN playback_target_url TEXT NOT NULL DEFAULT ''"); err != nil {
			return err
		}
	}

	var hasPlaybackModeColumn int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sites') WHERE name='playback_mode'").Scan(&hasPlaybackModeColumn); err != nil {
		return err
	}
	if hasPlaybackModeColumn == 0 {
		if _, err := d.db.Exec("ALTER TABLE sites ADD COLUMN playback_mode TEXT NOT NULL DEFAULT 'direct'"); err != nil {
			return err
		}
	}

	var hasStreamHostsColumn int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('sites') WHERE name='stream_hosts'").Scan(&hasStreamHostsColumn); err != nil {
		return err
	}
	if hasStreamHostsColumn == 0 {
		if _, err := d.db.Exec("ALTER TABLE sites ADD COLUMN stream_hosts TEXT NOT NULL DEFAULT '[]'"); err != nil {
			return err
		}
	}

	var hasHourlyIndex int
	if err := d.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_traffic_site_hour'").Scan(&hasHourlyIndex); err != nil {
		return err
	}
	if hasHourlyIndex > 0 {
		return nil
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		CREATE TABLE traffic_logs_dedup (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			site_id INTEGER NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
			bytes_in BIGINT DEFAULT 0,
			bytes_out BIGINT DEFAULT 0,
			recorded_at DATETIME NOT NULL
		);
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO traffic_logs_dedup (site_id, bytes_in, bytes_out, recorded_at)
		SELECT site_id, SUM(bytes_in), SUM(bytes_out), recorded_at
		FROM traffic_logs
		GROUP BY site_id, recorded_at;
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		DELETE FROM traffic_logs;
		INSERT INTO traffic_logs (site_id, bytes_in, bytes_out, recorded_at)
		SELECT site_id, bytes_in, bytes_out, recorded_at
		FROM traffic_logs_dedup;
		DROP TABLE traffic_logs_dedup;
		CREATE UNIQUE INDEX idx_traffic_site_hour ON traffic_logs(site_id, recorded_at);
	`); err != nil {
		return err
	}

	return tx.Commit()
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// Models
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

type TrafficLog struct {
	ID         int64  `json:"id"`
	SiteID     int64  `json:"site_id"`
	BytesIn    int64  `json:"bytes_in"`
	BytesOut   int64  `json:"bytes_out"`
	RecordedAt string `json:"recorded_at"`
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// DB operations
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

func (d *DB) AddTraffic(siteID, bytesIn, bytesOut int64) {
	if err := d.addTraffic(siteID, bytesIn, bytesOut); err != nil {
		log.Printf("[traffic] failed to persist usage for site %d: %v", siteID, err)
	}
}

func (d *DB) addTraffic(siteID, bytesIn, bytesOut int64) error {
	hour := time.Now().Truncate(time.Hour).Format("2006-01-02 15:04:05")
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`INSERT INTO traffic_logs (site_id, bytes_in, bytes_out, recorded_at)
		 VALUES (?,?,?,?)
		 ON CONFLICT(site_id, recorded_at) DO UPDATE SET
		 	bytes_in = traffic_logs.bytes_in + excluded.bytes_in,
		 	bytes_out = traffic_logs.bytes_out + excluded.bytes_out`,
		siteID, bytesIn, bytesOut, hour,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		"UPDATE sites SET traffic_used=traffic_used+?+?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		bytesIn, bytesOut, siteID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DB) GetTrafficLogs(siteID int64, hours int) ([]TrafficLog, error) {
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Format("2006-01-02 15:04:05")
	rows, err := d.db.Query(
		"SELECT id, site_id, bytes_in, bytes_out, recorded_at FROM traffic_logs WHERE site_id=? AND recorded_at>=? ORDER BY recorded_at",
		siteID, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []TrafficLog
	for rows.Next() {
		var l TrafficLog
		rows.Scan(&l.ID, &l.SiteID, &l.BytesIn, &l.BytesOut, &l.RecordedAt)
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []TrafficLog{}
	}
	return logs, nil
}

func (d *DB) DashboardStats() map[string]interface{} {
	var total, online int
	d.db.QueryRow("SELECT COUNT(*) FROM sites").Scan(&total)
	d.db.QueryRow("SELECT COUNT(*) FROM sites WHERE enabled=1").Scan(&online)
	var totalTraffic int64
	d.db.QueryRow("SELECT COALESCE(SUM(traffic_used),0) FROM sites").Scan(&totalTraffic)
	return map[string]interface{}{
		"total_sites":   total,
		"online_sites":  online,
		"total_traffic": totalTraffic,
	}
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// Proxy Engine
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// Diagnostics
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?
type DiagResult struct {
	Upstreams DiagUpstreams `json:"upstreams"`
	Health    DiagHealth    `json:"health"`
	TLS       DiagTLS       `json:"tls"`
	Headers   DiagHeaders   `json:"headers"`
	Proxy     DiagProxy     `json:"proxy"`
}

type DiagUpstreams struct {
	Primary  DiagUpstream `json:"primary"`
	Playback DiagUpstream `json:"playback"`
}

type DiagUpstream struct {
	Configured    bool       `json:"configured"`
	ConfiguredURL string     `json:"configured_url,omitempty"`
	EffectiveURL  string     `json:"effective_url"`
	UsingFallback bool       `json:"using_fallback"`
	SameAsPrimary bool       `json:"same_as_primary"`
	ShowHealth    bool       `json:"show_health"`
	ShowTLS       bool       `json:"show_tls"`
	Health        DiagHealth `json:"health"`
	TLS           DiagTLS    `json:"tls"`
}

type DiagProbe struct {
	Kind       string `json:"kind"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	HTTPStatus int    `json:"http_status,omitempty"`
}

type DiagHealth struct {
	Status    string    `json:"status"` // online, offline, error
	EmbyVer   string    `json:"emby_version"`
	LatencyMs int64     `json:"latency_ms"`
	Probe     DiagProbe `json:"probe"`
	Error     string    `json:"error,omitempty"`
}

type DiagTLS struct {
	Enabled   bool   `json:"enabled"`
	Valid     bool   `json:"valid"`
	Issuer    string `json:"issuer"`
	ExpiresAt string `json:"expires_at"`
	DaysLeft  int    `json:"days_left"`
	Error     string `json:"error,omitempty"`
}

type DiagHeaders struct {
	UAApplied    bool   `json:"ua_applied"`
	CurrentUA    string `json:"current_ua"`
	ClientField  string `json:"client_field"`
	VersionField string `json:"version_field"`
}

type DiagProxy struct {
	Running    bool   `json:"running"`
	ListenPort int    `json:"listen_port"`
	TotalReqs  int64  `json:"total_requests"`
	Uptime     string `json:"uptime,omitempty"`
}

func tlsIssuerName(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}
	if len(cert.Issuer.Organization) > 0 && cert.Issuer.Organization[0] != "" {
		return cert.Issuer.Organization[0]
	}
	if cert.Issuer.CommonName != "" {
		return cert.Issuer.CommonName
	}
	return cert.Issuer.String()
}

func canonicalTargetKey(target *url.URL) string {
	if target == nil {
		return ""
	}

	normalized := *target
	normalized.Scheme = strings.ToLower(normalized.Scheme)
	normalized.Host = strings.ToLower(normalized.Host)
	normalized.RawQuery = ""
	normalized.Fragment = ""

	cleanPath := path.Clean("/" + strings.Trim(normalized.Path, "/"))
	if cleanPath == "." || cleanPath == "/" {
		normalized.Path = ""
	} else {
		normalized.Path = cleanPath
	}

	return normalized.String()
}

func buildProbeURLs(target *url.URL, suffixes []string) []string {
	basePath := strings.TrimSpace(target.Path)
	seen := map[string]struct{}{}
	urls := make([]string, 0, len(suffixes))
	for _, suffix := range suffixes {
		probe := *target
		probe.RawQuery = ""
		probe.Fragment = ""
		if suffix == "" {
			cleanPath := path.Clean("/" + strings.Trim(basePath, "/"))
			if cleanPath == "." || cleanPath == "" {
				cleanPath = "/"
			}
			probe.Path = cleanPath
		} else {
			probe.Path = path.Clean("/" + path.Join(strings.Trim(basePath, "/"), suffix))
		}
		if _, ok := seen[probe.String()]; ok {
			continue
		}
		seen[probe.String()] = struct{}{}
		urls = append(urls, probe.String())
	}
	return urls
}

func healthProbeURLs(target *url.URL) []string {
	if strings.TrimSpace(target.Path) == "" || strings.TrimSpace(target.Path) == "/" {
		return buildProbeURLs(target, []string{"System/Info/Public", "emby/System/Info/Public", ""})
	}
	return buildProbeURLs(target, []string{"System/Info/Public", ""})
}

func playbackProbeURLs(target *url.URL) []string {
	return healthProbeURLs(target)
}

type diagProbePlan struct {
	BaseURL       string
	Kind          string
	Method        string
	CandidateURLs []string
	ParseVersion  bool
}

func resolveProbeKind(plan diagProbePlan, probeURL string) string {
	if plan.Kind != "metadata_api" {
		return plan.Kind
	}

	baseTarget, baseErr := normalizeTargetURL(plan.BaseURL)
	probeTarget, probeErr := normalizeTargetURL(probeURL)
	if baseErr != nil || probeErr != nil {
		return plan.Kind
	}

	basePath := strings.TrimSpace(baseTarget.Path)
	if basePath == "" {
		basePath = "/"
	}
	probePath := strings.TrimSpace(probeTarget.Path)
	if probePath == "" {
		probePath = "/"
	}
	if strings.TrimRight(probePath, "/") == strings.TrimRight(basePath, "/") {
		return "reachability_fallback"
	}

	return plan.Kind
}

func probeStatusRank(status int) int {
	switch {
	case status >= 200 && status < 300:
		return 4
	case status == http.StatusUnauthorized || status == http.StatusForbidden || status == http.StatusMethodNotAllowed:
		return 3
	case status == http.StatusNotFound:
		return 2
	case status > 0 && status < 500:
		return 1
	default:
		return 0
	}
}

func probeTargetHealth(plan diagProbePlan) DiagHealth {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	var bestReachable DiagHealth
	bestReachableRank := 0
	var serverError DiagHealth

	for _, probeURL := range plan.CandidateURLs {
		health := DiagHealth{
			Probe: DiagProbe{
				Kind:   resolveProbeKind(plan, probeURL),
				Method: plan.Method,
				URL:    probeURL,
			},
		}
		req, err := http.NewRequest(plan.Method, probeURL, nil)
		if err != nil {
			health.Status = "offline"
			health.Error = err.Error()
			return health
		}

		start := time.Now()
		resp, err := client.Do(req)
		latency := time.Since(start).Milliseconds()
		health.LatencyMs = latency
		if err != nil {
			health.Status = "offline"
			health.Error = err.Error()
			return health
		}

		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		health.Probe.HTTPStatus = resp.StatusCode

		if resp.StatusCode >= 500 {
			if serverError.Error == "" {
				health.Status = "error"
				health.Error = fmt.Sprintf("probe returned HTTP %d", resp.StatusCode)
				serverError = health
			}
			continue
		}

		health.Status = "online"
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if plan.ParseVersion {
				var info map[string]interface{}
				if json.Unmarshal(body, &info) == nil {
					if v, ok := info["Version"]; ok {
						health.EmbyVer = fmt.Sprintf("%v", v)
					}
				}
			}
			return health
		}

		rank := probeStatusRank(resp.StatusCode)
		if rank > bestReachableRank {
			bestReachable = health
			bestReachableRank = rank
		}
		if plan.Kind == "playback_path" && rank >= 3 {
			return health
		}
	}

	if bestReachableRank > 0 {
		return bestReachable
	}
	if serverError.Error != "" {
		return serverError
	}
	return DiagHealth{
		Status: "offline",
		Probe: DiagProbe{
			Kind:   plan.Kind,
			Method: plan.Method,
			URL:    plan.BaseURL,
		},
		Error: "health probe failed",
	}
}

func probeSiteHealth(targetURL string) DiagHealth {
	target, err := normalizeTargetURL(targetURL)
	if err != nil {
		return DiagHealth{
			Status: "offline",
			Probe: DiagProbe{
				Kind:   "metadata_api",
				Method: http.MethodGet,
			},
			Error: err.Error(),
		}
	}
	return probeTargetHealth(diagProbePlan{
		BaseURL:       target.String(),
		Kind:          "metadata_api",
		Method:        http.MethodGet,
		CandidateURLs: healthProbeURLs(target),
		ParseVersion:  true,
	})
}

func probePlaybackHealth(targetURL string) DiagHealth {
	target, err := normalizeTargetURL(targetURL)
	if err != nil {
		return DiagHealth{
			Status: "offline",
			Probe: DiagProbe{
				Kind:   "metadata_api",
				Method: http.MethodGet,
			},
			Error: err.Error(),
		}
	}
	return probeTargetHealth(diagProbePlan{
		BaseURL:       target.String(),
		Kind:          "metadata_api",
		Method:        http.MethodGet,
		CandidateURLs: playbackProbeURLs(target),
		ParseVersion:  true,
	})
}

func probeSiteTLS(target *url.URL) DiagTLS {
	var result DiagTLS
	if target == nil || !strings.EqualFold(target.Scheme, "https") {
		return result
	}

	result.Enabled = true
	host := target.Hostname()
	port := target.Port()
	if port == "" {
		port = "443"
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", host+":"+port, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return result
	}

	cert := certs[0]
	result.Valid = time.Now().Before(cert.NotAfter)
	result.Issuer = tlsIssuerName(cert)
	result.ExpiresAt = cert.NotAfter.Format("2006-01-02")
	result.DaysLeft = int(time.Until(cert.NotAfter).Hours() / 24)

	return result
}

func diagnoseUpstreamTarget(targetURL, probeKind string) (DiagUpstream, string) {
	trimmed := strings.TrimSpace(targetURL)
	result := DiagUpstream{
		Configured:    trimmed != "",
		ConfiguredURL: trimmed,
		EffectiveURL:  trimmed,
		ShowHealth:    true,
	}

	parsed, err := normalizeTargetURL(targetURL)
	if err != nil {
		result.Health = DiagHealth{Status: "offline", Error: err.Error()}
		return result, ""
	}

	result.ConfiguredURL = parsed.String()
	result.EffectiveURL = parsed.String()
	switch probeKind {
	case "playback_path":
		result.Health = probePlaybackHealth(parsed.String())
	default:
		result.Health = probeSiteHealth(parsed.String())
	}
	result.TLS = probeSiteTLS(parsed)
	result.ShowTLS = result.TLS.Enabled

	return result, canonicalTargetKey(parsed)
}

func diagnoseSite(site *Site, pm *ProxyManager) DiagResult {
	profile := getUAProfile(site.UAMode)
	primary, primaryKey := diagnoseUpstreamTarget(site.TargetURL, "metadata_api")
	primary.Configured = true
	primary.ShowHealth = true
	primary.ShowTLS = primary.TLS.Enabled

	playbackRaw := strings.TrimSpace(site.PlaybackTargetURL)
	playback := primary
	playback.ConfiguredURL = ""
	playback.Configured = false
	playback.UsingFallback = true
	playback.SameAsPrimary = true
	playback.ShowHealth = false
	playback.ShowTLS = false

	if playbackRaw != "" {
		var playbackKey string
		playback, playbackKey = diagnoseUpstreamTarget(playbackRaw, "playback_path")
		playback.Configured = true
		playback.UsingFallback = false
		playback.SameAsPrimary = playbackKey != "" && playbackKey == primaryKey
		if playback.SameAsPrimary {
			playback.Health = primary.Health
			playback.TLS = primary.TLS
			playback.EffectiveURL = primary.EffectiveURL
			playback.ShowHealth = false
			playback.ShowTLS = false
		}
	}

	result := DiagResult{
		Upstreams: DiagUpstreams{
			Primary:  primary,
			Playback: playback,
		},
		Health: primary.Health,
		TLS:    primary.TLS,
	}

	// Headers
	result.Headers = DiagHeaders{
		UAApplied:    true,
		CurrentUA:    profile.UserAgent,
		ClientField:  profile.Client,
		VersionField: profile.Version,
	}

	// Proxy status
	result.Proxy = DiagProxy{
		Running:    pm.IsRunning(site.ID),
		ListenPort: site.ListenPort,
	}

	return result
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// HTTP API
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?

func originMatchesRequestHost(origin string, r *http.Request) bool {
	parsed, err := url.Parse(origin)
	if err != nil || parsed.User != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	if parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	return strings.EqualFold(parsed.Host, r.Host)
}

func refererMatchesRequestHost(referer string, r *http.Request) bool {
	parsed, err := url.Parse(referer)
	if err != nil || parsed.User != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	return strings.EqualFold(parsed.Host, r.Host)
}

func requestHasSameOrigin(r *http.Request) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		return originMatchesRequestHost(origin, r)
	}
	return refererMatchesRequestHost(r.Referer(), r)
}

func stateChangingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func panelCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if !originMatchesRequestHost(origin, r) {
				http.Error(w, "cross-origin request denied", http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", panelCSP)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		next.ServeHTTP(w, r)
	})
}

func requestIsHTTPS(r *http.Request) bool {
	if r.TLS != nil || strings.EqualFold(r.URL.Scheme, "https") {
		return true
	}
	return false
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("empty body")
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	return json.NewDecoder(r.Body).Decode(dst)
}

type App struct {
	db           *DB
	pm           *ProxyManager
	loginLimiter *loginRateLimiter
	setupToken   string
	setupTokenMu sync.Mutex
}

func (a *App) jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (a *App) jsonErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// GET /api/dashboard
func (a *App) handleDashboard(w http.ResponseWriter, r *http.Request) {
	stats := a.db.DashboardStats()
	stats["running_sites"] = a.pm.GetRunningCount()
	a.jsonOK(w, stats)
}

// ExportSiteRecord is the JSON structure for backup/restore

// GET /api/traffic/{site_id}
func (a *App) handleTraffic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/traffic/")

	if path == "overview" {
		stats := a.db.DashboardStats()
		a.jsonOK(w, stats)
		return
	}

	siteID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		a.jsonErr(w, 400, "invalid site id")
		return
	}

	hours := 24
	if h := r.URL.Query().Get("hours"); h != "" {
		if v, err := strconv.Atoi(h); err == nil {
			hours = v
		}
	}

	logs, err := a.db.GetTrafficLogs(siteID, hours)
	if err != nil {
		a.jsonErr(w, 500, err.Error())
		return
	}
	a.jsonOK(w, logs)
}

// GET /api/ua-profiles
func (a *App) handleUAProfiles(w http.ResponseWriter, r *http.Request) {
	profiles := make([]UAProfile, 0, len(uaProfiles))
	for _, p := range uaProfiles {
		profiles = append(profiles, p)
	}
	a.jsonOK(w, profiles)
}

// GET /api/events — Server-Sent Events stream
func (a *App) handleSSE(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		a.jsonErr(w, 405, "method not allowed")
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		a.jsonErr(w, 500, "SSE not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()

	// Send initial data immediately
	a.sendSSEEvent(w, flusher)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.sendSSEEvent(w, flusher)
		}
	}
}

func (a *App) sendSSEEvent(w http.ResponseWriter, flusher http.Flusher) {
	stats := a.db.DashboardStats()
	stats["running_sites"] = a.pm.GetRunningCount()
	stats["total_requests"] = a.pm.GetTotalRequests()
	stats["uptime_seconds"] = int(time.Since(startTime).Seconds())

	// Collect per-site live stats
	a.pm.mu.RLock()
	siteStats := make([]map[string]interface{}, 0)
	for _, inst := range a.pm.proxies {
		siteStats = append(siteStats, map[string]interface{}{
			"id":        inst.Site.ID,
			"name":      inst.Site.Name,
			"bytes_in":  inst.bytesIn.Load(),
			"bytes_out": inst.bytesOut.Load(),
			"requests":  inst.reqCount.Load(),
			"running":   true,
		})
	}
	a.pm.mu.RUnlock()
	stats["live_sites"] = siteStats

	data, _ := json.Marshal(stats)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}

// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?// Main 闂?with graceful shutdown
// 闂傚倸鍊搁崐椋庣矆娴ｅ搫顥氭い鎾卞灩绾惧潡鏌曢崼婵愭Ц缂佲偓婢舵劗鍙撻柛銉ｅ妿閳藉鏌ｉ妶澶岀暫闁哄矉绱曟禒锔炬嫚閹绘帒顫撻梻浣虹帛閹稿鎯勯鐐茶摕闁绘柨鍚嬮崵瀣亜閹哄棗浜炬繝寰枫倕袚缂佺粯鐩畷銊╊敊閸撗呭帨闂備礁鎼懟顖滅矓瑜版帒绠栨繝濠傚悩閻旂厧浼犻柛鏇炵仛缂嶅倿姊婚崒娆戭槮闁圭⒈鍋婇獮濠呯疀濞戞瑥浜楅梺璺ㄥ枔婵挳寮伴妷鈺傜叆闁绘柨鎼瓭缂備胶濮甸惄顖炲蓟閺囩喓绡€闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及閵夆晜鐓ラ柣鏂挎惈瀛濈紓浣哄У閻╊垶寮婚弴鐔虹瘈闊洦绋掗宥夋⒑缂佹绠栧┑鐐诧工椤繘宕崟顓熸闂佹悶鍎滈崘顭戠€遍梻鍌欑閹诧繝寮婚妸褎宕叉俊顖欒閸ゆ洟鏌＄仦璇插姎闁藉啰鍠栭弻鏇熷緞閸繂濮㈤梺鍛娚戦幃鍌氼潖閾忚鍏滈柛娑卞幘閸旂兘姊洪崨濠冪叆缂佸鎸抽崺銏狀吋閸滀胶鍙嗛梺鍓插亞閸犳捇宕㈤幘缁樷拺缂備焦锚閻忥箓鏌ㄥ顑芥斀妞ゆ梻鎳撴禍楣冩⒒閸屾瑧顦﹂柟纰卞亰楠炲﹨绠涘☉娆忎簵闂佽法鍠撴慨鎾及?
var startTime = time.Now()

func runCommandLine(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	switch args[0] {
	case "--version", "-v":
		if len(args) != 1 {
			return true, errors.New("version command does not accept arguments")
		}
		fmt.Println(appVersion)
		return true, nil
	case "--healthcheck":
		if len(args) != 1 {
			return true, errors.New("healthcheck command does not accept arguments")
		}
		return true, runHealthcheckCommand()
	default:
		return false, nil
	}
}

func runHealthcheckCommand() error {
	port := 9090
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
	client := &http.Client{
		Timeout: 4 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var lastErr error
	for _, host := range []string{"127.0.0.1", "::1"} {
		endpoint := fmt.Sprintf("http://%s/api/auth/check", net.JoinHostPort(host, strconv.Itoa(port)))
		resp, err := client.Get(endpoint)
		if err != nil {
			lastErr = err
			continue
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		closeErr := resp.Body.Close()
		if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices && closeErr == nil {
			return nil
		}
		lastErr = fmt.Errorf("%s returned HTTP %d", endpoint, resp.StatusCode)
	}
	if lastErr == nil {
		lastErr = errors.New("no loopback endpoints probed")
	}
	return fmt.Errorf("panel healthcheck failed: %w", lastErr)
}

func main() {
	if handled, err := runCommandLine(os.Args[1:]); handled {
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	port := 9090
	dbPath := resolveDefaultDBPath()
	if jwtSecretEphemeral {
		log.Printf("JWT_SECRET not set; generated an ephemeral signing secret for this process. Set JWT_SECRET explicitly for stable sessions.")
	}

	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		dbPath = v
	}

	for i, arg := range os.Args[1:] {
		switch arg {
		case "--port", "-p":
			if i+2 < len(os.Args) {
				if p, err := strconv.Atoi(os.Args[i+2]); err == nil {
					port = p
				}
			}
		case "--db":
			if i+2 < len(os.Args) {
				dbPath = os.Args[i+2]
			}
		}
	}

	bindAddr := strings.TrimSpace(os.Getenv("PANEL_BIND_ADDR"))
	addr, err := panelListenAddr(bindAddr, port)
	if err != nil {
		log.Fatalf("invalid panel listen address: %v", err)
	}
	if !panelBindIsLoopback(bindAddr) && !envBool("ALLOW_INSECURE_HTTP") {
		log.Fatalf("refusing to expose the management panel over plain HTTP on %s; set ALLOW_INSECURE_HTTP=true only for a deliberate temporary compatibility deployment, or bind PANEL_BIND_ADDR=127.0.0.1 behind an HTTPS reverse proxy", bindAddr)
	}

	trustedNets, err := parseTrustedProxyCIDRs(os.Getenv("TRUSTED_PROXY_CIDRS"))
	if err != nil {
		log.Fatalf("invalid TRUSTED_PROXY_CIDRS: %v", err)
	}
	trustedProxyNets = trustedNets
	if len(trustedProxyNets) > 0 {
		log.Printf("TRUSTED_PROXY_CIDRS: trusting %d CIDR(s) for X-Real-IP / X-Forwarded-For on login rate limits", len(trustedProxyNets))
	}

	db, err := openDB(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	setupToken, setupGenerated, err := resolveSetupToken(db.UserCount(), os.Getenv("SETUP_TOKEN"))
	if err != nil {
		log.Fatalf("initial setup unavailable: %v", err)
	}
	if setupGenerated {
		log.Printf("SETUP_TOKEN not set; generated %s — pass this token when creating the first administrator", setupToken)
	}

	pm := NewProxyManager(db)
	pm.StartAllEnabled()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				pm.FlushTraffic()
			case <-ctx.Done():
				return
			}
		}
	}()

	app := &App{db: db, pm: pm, loginLimiter: newLoginRateLimiter(), setupToken: setupToken}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/setup", panelCORS(app.handleSetup))
	mux.HandleFunc("/api/auth/login", panelCORS(app.handleLogin))
	mux.HandleFunc("/api/auth/logout", panelCORS(app.handleLogout))
	mux.HandleFunc("/api/auth/check", panelCORS(app.handleAuthCheck))

	mux.HandleFunc("/api/dashboard", panelCORS(app.authMiddleware(app.handleDashboard)))
	mux.HandleFunc("/api/sites", panelCORS(app.authMiddleware(app.handleSites)))
	mux.HandleFunc("/api/sites/export", panelCORS(app.authMiddleware(app.handleSitesExport)))
	mux.HandleFunc("/api/sites/import", panelCORS(app.authMiddleware(app.handleSitesImport)))
	mux.HandleFunc("/api/sites/", panelCORS(app.authMiddleware(app.handleSiteByID)))
	mux.HandleFunc("/api/traffic/", panelCORS(app.authMiddleware(app.handleTraffic)))
	mux.HandleFunc("/api/ua-profiles", panelCORS(app.authMiddleware(app.handleUAProfiles)))
	mux.HandleFunc("/api/events", panelCORS(app.authMiddleware(app.handleSSE)))

	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatalf("failed to load embedded static files: %v", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := staticFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       120 * time.Second,
	}

	log.Println("============================================================")
	log.Printf("  StreamDock - Emby reverse proxy management panel %s", appVersion)
	log.Printf("  Listening on: http://%s", addr)
	log.Printf("  Sites loaded: %d (%d running)", func() int { s, _ := db.ListSites(); return len(s) }(), pm.GetRunningCount())
	log.Println("  Features: WebSocket proxy, TLS diagnostics, traffic limits")
	log.Println("============================================================")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-sigCh
	log.Println("\nReceived shutdown signal, stopping StreamDock...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	pm.GracefulShutdown(shutdownCtx)
	srv.Shutdown(shutdownCtx)

	log.Println("StreamDock stopped cleanly")
}

func resolveDefaultDBPath() string {
	if _, err := os.Stat(defaultDBFile); err == nil {
		return defaultDBFile
	}
	if _, err := os.Stat(legacyDBFile); err == nil {
		log.Printf("using legacy database %s; new default is %s", legacyDBFile, defaultDBFile)
		return legacyDBFile
	}
	return defaultDBFile
}
