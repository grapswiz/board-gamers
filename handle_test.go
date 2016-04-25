package board_gamers

import (
	"reflect"
	"testing"
)

func TestExtractTrickplayGames(t *testing.T) {
	want := []string{"HAWAIIミニ拡張", "ロシアンレールロードミニ拡張＆ストーンエイジミニ拡張", "ヘックメック拡張"}
	text := "新しい神タイルや島タイルが含まれる「HAWAIIミニ拡張」、新しい技術者とボーナスタイルのセット「ロシアンレールロードミニ拡張＆ストーンエイジミニ拡張」、「ヘックメック拡張」が入荷しております。よろしくお願い致します。"
	if result := extractTrickplayGames(text); !reflect.DeepEqual(result, want) {
		t.Errorf("extractTrickplayGames = %v, want %v", result, want)
	}

	want = []string{"T.I.M.E Stories", "T.I.M.E Storiesシナリオ The Marcy Case"}
	text = "Space Cowboysが贈る壮大な謎解きゲーム「T.I.M.E Stories」、今度は異なる世界の過去の地球において、マーシィという女性を救う「T.I.M.E Storiesシナリオ The Marcy Case」 #トリックプレイ"
	if result := extractTrickplayGames(text); !reflect.DeepEqual(result, want) {
		t.Errorf("extractTrickplayGames = %v, want %v", result, want)
	}
}

func TestExtractTendaysGames(t *testing.T) {
	want := []string{"ゲームマーケットカタログ", "ゲームマーケットホールマップ", "ナショナルエコノミー", "ドミニオンマニアックスSpecial"}
	text := "国内最大級のボードゲームイベントまでもうすぐ！「ゲームマーケットカタログ」、「ゲームマーケットホールマップ」を新入荷しました。どちらも入場券を兼ねています。\nナショナルエコノミー、ドミニオンマニアックスSpecialを再入荷しました。"
	if result := extractTendaysGames(text); !reflect.DeepEqual(result, want) {
		t.Errorf("extractTendaysGames = %v, want %v", result, want)
	}

	want = []string{"大いなる狂気の書日本語版", "スチームタイム", "二枚目が好き", "山頂をめざせ", "おしくらモンスター", "双天至尊堂・天九牌", "カルカソンヌ", "お邪魔者", "8か28", "ワードバスケット"}
	text = "デッキ構築し、協力して迫りくる魔物たちを撃退しろ！「大いなる狂気の書日本語版」、「スチームタイム」、「二枚目が好き」、「山頂をめざせ」、「おしくらモンスター」、「双天至尊堂・天九牌」を新入荷しました。\nカルカソンヌ、お邪魔者、8か28、ワードバスケットを再入荷しました。"
	if result := extractTendaysGames(text); !reflect.DeepEqual(result, want) {
		t.Errorf("extractTendaysGames = %v, want %v", result, want)
	}
}
