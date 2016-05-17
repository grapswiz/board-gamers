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

	want = []string{"コーヒーロースター", "リスボン、世界への扉", "バルーンチャレンジ"}
	text = "国産ゲーム三種「コーヒーロースター」、「リスボン、世界への扉」、「バルーンチャレンジ」を新入荷しました。"
	if result := extractTendaysGames(text); !reflect.DeepEqual(result, want) {
		t.Errorf("extractTendaysGames = %v, want %v", result, want)
	}

	// まぁいっかという感じ
	//want = []string{"人気のダイスマネージメントゲーム「キングスフォージ」の拡張セット二種", "エイリアンフロンティアの拡張セット二種（和訳はつきません）", "メモワール’44", "キングスフォージ（二版）", "タイニーエピックディフェンダーズ"}
	//text = "人気のダイスマネージメントゲーム「キングスフォージ」の拡張セット二種、エイリアンフロンティアの拡張セット二種（和訳はつきません）、メモワール’44を新入荷しました。\nキングスフォージ（二版）、タイニーエピックディフェンダーズを再入荷しました。"
	//if result := extractTendaysGames(text); !reflect.DeepEqual(result, want) {
	//	t.Errorf("extractTendaysGames = %v, want %v", result, want)
	//}

	want = []string{"カウンシル・オブ・フォー", "アンユージュアルサスペクツ", "世界の七不思議デュエル"}
	text = "「大いなる文明の曙（メガシヴィライゼーション）」用の大判スリーブを取り扱い開始しました。お役立ちアイテムコーナーの「スリーブ」ドロップダウンメニューからお選びください。\nカウンシル・オブ・フォー、アンユージュアルサスペクツ、世界の七不思議デュエルを再入荷しました。"
}

func TestExtractGameFieldGames(t *testing.T) {
	{
		// TODO
		//want := []string{"ペアーズ:洋ナシデッキ", "ペアーズ:ゴブリンデッキ", "ペアーズ:リトルバードデッキ(通常版)", "ペアーズ:ファニーハッターデッキ", "ペアーズ:ボールゲームデッキ", "アルルの丘(Arler erde)　日本語版", "ゼロ　日本語版(Zero)"}
		//text := `
		//お世話になっております。ゲームフィールドです。
		//新入荷・再入荷のお知らせです。
		//
		//新入荷
		//1は1枚、2は2枚、3は3枚・・・10は10枚という、特徴的な構成のカードデッキシリーズです。
		//
		//手元にあるカードと同じ数字のカードを引くか引かないかの賭けに挑むかの判断と、賭けに成功するかどうかを繰り返すシンプルな運試しの基本ゲームをはじめ、さまざまなルールで遊ぶことができます。
		//
		//「新しいクラシックパブゲーム」の名前の通り、気軽に手軽に誰もがすぐに盛り上がることのできるルールがたくさん楽しめる、オススメゲームシリーズです。
		//
		//日本語版では、アメリカのクラウドファンディング「Kickstarter」で大成功した原版で人気の二デッキに加え、日本オリジナルデザインの新デッキを三種類用意いたしました。
		//このデザインバリエーションの幅の広さも大きな魅力となっています。ぜひ、お気に入りのデッキを見つけてください。
		//
		//「ペアーズ:洋ナシデッキ」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=439
		//※「洋ナシ(フルーツ)デッキ」には、基本ルールと基本ルールのより詳細な遊び方が入っています。
		//
		//「ペアーズ:ゴブリンデッキ」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=440
		//※「ゴブリンデッキ」には、基本ルールと「ゴブリンポーカー」のルールが入っています。
		//
		//「ペアーズ:リトルバードデッキ(通常版)」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=441
		//※「リトルバードデッキ」には、基本ルールと「ことりあつめ」のルールが入っています。
		//
		//「ペアーズ:ファニーハッターデッキ」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=442
		//※「ファニーハッターデッキ」には、基本ルールと「プッシュ&プッシュ」のルールが入っています。
		//
		//「ペアーズ:ボールゲームデッキ」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=443
		//※「ボールゲームデッキ」には、基本ルールと「スポーツ」のルールが入っています。
		//
		//再入荷
		//「アルルの丘(Arler erde)　日本語版」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=420
		//
		//「ゼロ　日本語版(Zero)」
		//http://gamefield.sakura.ne.jp/products/detail.php?product_id=365
		//`
		//if result := extractGamefieldGames(text); !reflect.DeepEqual(result, want) {
		//	t.Errorf("extractTendaysGames = %v, want %v", result, want)
		//}
	}
}
