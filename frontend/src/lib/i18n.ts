import { useState, useEffect } from 'react'
import type { LobbyPhase } from '../types/lobby'

export type Language = 'ko' | 'ja' | 'pl' | 'zh' | 'en'

export const SUPPORTED_LANGUAGES: Record<Language, string> = {
  en: 'English',
  ko: '한국어',
  ja: '日本語',
  pl: 'Polski',
  zh: '中文',
}

const STORAGE_KEY = 'croquing_language'
const listeners = new Set<() => void>()

export function getBrowserLanguage(): Language {
  const lang = navigator.language.toLowerCase()
  if (lang.startsWith('ko')) return 'ko'
  if (lang.startsWith('ja')) return 'ja'
  if (lang.startsWith('pl')) return 'pl'
  if (lang.startsWith('zh')) return 'zh'
  return 'en'
}

export function getCurrentLanguage(): Language {
  const saved = localStorage.getItem(STORAGE_KEY) as Language
  if (saved && SUPPORTED_LANGUAGES[saved]) {
    return saved
  }
  return getBrowserLanguage()
}

export function setLanguage(lang: Language): void {
  if (SUPPORTED_LANGUAGES[lang]) {
    localStorage.setItem(STORAGE_KEY, lang)
    syncDocumentLanguage()
    listeners.forEach((listener) => listener())
  }
}

export function useLanguage(): Language {
  const [lang, setLang] = useState<Language>(getCurrentLanguage())

  useEffect(() => {
    const handleUpdate = () => {
      setLang(getCurrentLanguage())
    }
    listeners.add(handleUpdate)
    return () => {
      listeners.delete(handleUpdate)
    }
  }, [])

  return lang
}

export function syncDocumentLanguage(): void {
  document.documentElement.lang = getCurrentLanguage()
}

// Simple translation lookup with placeholder support like {count}
const TRANSLATIONS: Record<Language, Record<string, string>> = {
  en: {
    // HomePage
    'home.eyebrow': 'Real-time croquis meetups',
    'home.lead':
      'Create a lobby, share the link, and draw together with synchronized photos and timers — no screen sharing required.',
    'home.createLobby': 'Create lobby',
    'home.createLobbyFailed': 'Failed to create lobby',
    'home.creatingLobby': 'Creating lobby…',
    'home.howItWorks': 'How it works',
    'home.step1.title': 'Create & share',
    'home.step1.desc': 'Start a lobby and send the link to your drawing group.',
    'home.step2.title': 'Pick references',
    'home.step2.desc':
      'The admin selects photos from Pixabay for everyone to draw.',
    'home.step3.title': 'Draw in sync',
    'home.step3.desc':
      'Timed rounds with the same photo and countdown for all.',

    // LobbyLayout & Page
    'lobby.badge.admin': 'Admin',
    'lobby.badge.participant': 'Participant',
    'lobby.connection.connecting': 'Connecting…',
    'lobby.connection.connected': 'Live',
    'lobby.connection.reconnecting': 'Reconnecting…',
    'lobby.connection.disconnected': 'Disconnected',
    'lobby.connection.lost': 'Connection lost. Trying to reconnect…',
    'lobby.loadingState': 'Loading lobby state…',
    'lobby.invalidLink': 'Invalid lobby link.',
    'lobby.backHome': 'Back home',
    'lobby.participantCount': '{count} participant',
    'lobby.participantCountPlural': '{count} participants',
    'lobby.lobbyId': 'Lobby {id}',

    // CopyLobbyLinkButton
    'copy.copied': 'Copied!',
    'copy.failed': 'Copy failed',
    'copy.link': 'Copy link',

    // AdminControls
    'admin.startSession': 'Start session',
    'admin.starting': 'Starting…',
    'admin.nextPhoto': 'Next photo',
    'admin.loading': 'Loading…',
    'admin.endSession': 'End session',
    'admin.ending': 'Ending…',
    'admin.actionFailed': 'Action failed',

    // RoundDurationPicker
    'duration.label': 'Round time',
    'duration.unit': 'min',
    'duration.minutes': '{count} min per round',
    'duration.ariaGroup': 'Round duration in minutes',
    'duration.decrease': 'Decrease round time',
    'duration.increase': 'Increase round time',
    'duration.hint': 'Use the left and right buttons to adjust.',
    'duration.updateFailed': 'Failed to update round time',

    // ParticipantWaitPanel
    'wait.settingPhotos': 'Setting up reference photos...',
    'wait.adminChoosing':
      'The admin is choosing photos. They stay hidden until the session starts.',
    'wait.tipLabel': 'Tip: {title}',
    'wait.goTip': 'Go to tip {idx}',
    'wait.waitingForAdmin': 'Waiting for admin...',
    'wait.carouselNav': 'Drawing tips',
    'wait.pauseCarousel': 'Pause tip rotation',
    'wait.resumeCarousel': 'Resume tip rotation',

    // PhotoReviewPanel
    'review.photosSaved': '{count} photos saved',
    'review.instruction':
      'Hover a thumbnail to preview it at full size. When you are happy with the set, complete the selection to prepare the session (photos will be shuffled automatically).',
    'review.previewAria': 'Preview saved photo {index} of {total}',
    'review.editSelection': 'Edit selection',
    'review.selectionComplete': 'Selection complete',
    'review.shuffling': 'Shuffling…',
    'review.modalAria': 'Photo preview',
    'review.closePreview': 'Close preview',
    'review.modalAlt': 'Reference photo {index} of {total}',

    // PhotoSelectionPanel
    'selection.errSave': 'Failed to save selection',
    'selection.errConfirm': 'Failed to confirm selection',
    'selection.errReopen': 'Failed to reopen photo selection',
    'selection.saving': 'Saving…',
    'selection.saveCount': 'Save {count} photos',

    // ReadyPanel
    'ready.photosReady': 'photos ready',
    'ready.desc':
      'The order is shuffled and hidden. Thumbnails stay off until each draw round begins.',
    'ready.hint': 'Waiting for the admin to start…',

    // DrawingPanel
    'draw.exitFullscreen': 'Exit Fullscreen',
    'draw.enterFullscreen': 'Enter Fullscreen',
    'draw.waitingPhoto': 'Waiting for photo…',
    'draw.attribution': 'Image from',
    'draw.round': 'Round {current} / {total}',
    'draw.startsIn': 'Round starts in {count} seconds',
    'draw.remainingAria': 'Draw time remaining',
    'draw.photoAlt': 'Reference photo for this croquis round',

    // SessionBreakPanels
    'break.takeBreather': 'Take a breather',
    'break.hiddenDesc':
      'The reference photo is hidden until the next round starts.',
    'break.completedRound': 'Completed round {current} of {total}',
    'break.sessionFinished': 'Session finished',
    'break.completedRoundsDesc':
      'You completed {count} round. Great work everyone.',
    'break.completedRoundsDescPlural':
      'You completed {count} rounds. Great work everyone.',
    'break.downloadPhotos': 'Download reference photos (ZIP)',

    // PixabaySearchPanel
    'search.errEmpty': 'Enter a search term',
    'search.errFailed': 'Search failed',
    'search.fieldLabel': 'Search Pixabay',
    'search.placeholder': 'e.g. portrait, anatomy, gesture',
    'search.sort': 'Sort',
    'search.sortPopular': 'Popular',
    'search.sortLatest': 'Latest',
    'search.searching': 'Searching',
    'search.button': 'Search',
    'search.prevPage': 'Previous page',
    'search.nextPage': 'Next page',
    'search.pageIndicator': 'Page {page} of {total}',
    'search.hint': '{count} selected · {recommended} recommended',
    'search.ariaSelect': 'Select photo {id}',
    'search.ariaDeselect': 'Deselect photo {id}',
    'search.dock.title': 'Selected Reference Photos',
    'search.dock.remove': 'Remove photo',
    'search.dock.removeAria': 'Remove photo {id}',
  },
  ko: {
    // HomePage
    'home.eyebrow': '실시간 크로키 모임',
    'home.lead':
      '로비를 만들고 링크만 공유하면, 같은 사진과 타이머로 함께 그릴 수 있어요. 화면 공유는 필요 없습니다.',
    'home.createLobby': '로비 만들기',
    'home.createLobbyFailed': '로비를 만들지 못했습니다',
    'home.creatingLobby': '로비 만드는 중…',
    'home.howItWorks': '이용 방법',
    'home.step1.title': '로비 만들고 공유하기',
    'home.step1.desc': '로비를 만들고 그림 모임에 링크를 보내세요.',
    'home.step2.title': '레퍼런스 선택',
    'home.step2.desc': '방장이 모두가 그릴 사진을 고릅니다.',
    'home.step3.title': '동시에 그리기',
    'home.step3.desc': '모두 같은 사진을 보며 카운트다운과 함께 그립니다.',

    // LobbyLayout & Page
    'lobby.badge.admin': '방장',
    'lobby.badge.participant': '참가자',
    'lobby.connection.connecting': '연결 중…',
    'lobby.connection.connected': '연결됨',
    'lobby.connection.reconnecting': '재연결 중…',
    'lobby.connection.disconnected': '연결 끊김',
    'lobby.connection.lost': '연결이 끊겼습니다. 다시 연결하는 중…',
    'lobby.loadingState': '로비 불러오는 중…',
    'lobby.invalidLink': '유효하지 않은 로비 링크입니다.',
    'lobby.backHome': '홈으로',
    'lobby.participantCount': '{count}명 참가 중',
    'lobby.participantCountPlural': '{count}명 참가 중',
    'lobby.lobbyId': '로비 {id}',

    // CopyLobbyLinkButton
    'copy.copied': '복사했어요!',
    'copy.failed': '복사하지 못했습니다',
    'copy.link': '링크 복사',

    // AdminControls
    'admin.startSession': '세션 시작',
    'admin.starting': '시작 중…',
    'admin.nextPhoto': '다음 사진',
    'admin.loading': '불러오는 중…',
    'admin.endSession': '세션 종료',
    'admin.ending': '종료 중…',
    'admin.actionFailed': '작업에 실패했습니다',

    // RoundDurationPicker
    'duration.label': '라운드 시간',
    'duration.unit': '분',
    'duration.minutes': '라운드당 {count}분',
    'duration.ariaGroup': '라운드 시간(분)',
    'duration.decrease': '라운드 시간 줄이기',
    'duration.increase': '라운드 시간 늘리기',
    'duration.hint': '좌우 버튼으로 조절할 수 있어요.',
    'duration.updateFailed': '라운드 시간을 바꾸지 못했습니다',

    // ParticipantWaitPanel
    'wait.settingPhotos': '레퍼런스 사진 준비 중…',
    'wait.adminChoosing':
      '방장이 사진을 고르는 중이에요. 세션이 시작되기 전까지는 보이지 않습니다.',
    'wait.tipLabel': '팁: {title}',
    'wait.goTip': '팁 {idx}으로 이동',
    'wait.waitingForAdmin': '방장을 기다리는 중…',
    'wait.carouselNav': '그리기 팁',
    'wait.pauseCarousel': '팁 넘김 일시정지',
    'wait.resumeCarousel': '팁 넘김 재개',

    // PhotoReviewPanel
    'review.photosSaved': '사진 {count}장 저장됨',
    'review.instruction':
      '썸네일에 마우스를 올리면 크게 볼 수 있어요. 마음에 들면 선택 완료를 눌러 세션을 준비하세요. (사진 순서는 자동으로 섞이고 고정됩니다.)',
    'review.previewAria': '저장된 사진 {index} / {total} 미리보기',
    'review.editSelection': '선택 수정',
    'review.selectionComplete': '선택 완료',
    'review.shuffling': '섞는 중…',
    'review.modalAria': '사진 미리보기',
    'review.closePreview': '미리보기 닫기',
    'review.modalAlt': '레퍼런스 사진 {index} / {total}',

    // PhotoSelectionPanel
    'selection.errSave': '선택을 저장하지 못했습니다',
    'selection.errConfirm': '선택을 확정하지 못했습니다',
    'selection.errReopen': '사진 선택을 다시 열지 못했습니다',
    'selection.saving': '저장 중…',
    'selection.saveCount': '사진 {count}장 저장',

    // ReadyPanel
    'ready.photosReady': '장 사진 준비됨',
    'ready.desc':
      '사진 순서를 섞어 숨겨 두었어요. 각 라운드가 시작되기 전까지 썸네일은 보이지 않습니다.',
    'ready.hint': '방장이 시작할 때까지 기다려 주세요…',

    // DrawingPanel
    'draw.exitFullscreen': '전체 화면 끄기',
    'draw.enterFullscreen': '전체 화면',
    'draw.waitingPhoto': '사진 기다리는 중…',
    'draw.attribution': '사진 출처:',
    'draw.round': '라운드 {current} / {total}',
    'draw.startsIn': '{count}초 뒤 라운드 시작',
    'draw.remainingAria': '남은 그리기 시간',
    'draw.photoAlt': '이번 크로키 라운드 레퍼런스 사진',

    // SessionBreakPanels
    'break.takeBreather': '잠깐 쉬어 가세요',
    'break.hiddenDesc':
      '다음 라운드가 시작될 때까지 레퍼런스 사진은 숨겨져 있어요.',
    'break.completedRound': '{total}라운드 중 {current}라운드 완료',
    'break.sessionFinished': '세션 종료',
    'break.completedRoundsDesc': '총 {count}라운드 끝! 모두 수고하셨어요.',
    'break.completedRoundsDescPlural':
      '총 {count}라운드 끝! 모두 수고하셨어요.',
    'break.downloadPhotos': '레퍼런스 사진 ZIP 다운로드',

    // PixabaySearchPanel
    'search.errEmpty': '검색어를 입력하세요',
    'search.errFailed': '검색하지 못했습니다',
    'search.fieldLabel': 'Pixabay 검색',
    'search.placeholder': '예: 초상화, 해부학, 제스처',
    'search.sort': '정렬',
    'search.sortPopular': '인기순',
    'search.sortLatest': '최신순',
    'search.searching': '검색 중',
    'search.button': '검색',
    'search.prevPage': '이전 페이지',
    'search.nextPage': '다음 페이지',
    'search.pageIndicator': '{total}페이지 중 {page}페이지',
    'search.hint': '{count}장 선택 · {recommended}장 권장',
    'search.ariaSelect': '사진 {id} 선택',
    'search.ariaDeselect': '사진 {id} 선택 해제',
    'search.dock.title': '선택한 레퍼런스 사진',
    'search.dock.remove': '사진 빼기',
    'search.dock.removeAria': '사진 {id} 빼기',
  },
  ja: {
    // HomePage
    'home.eyebrow': 'リアルタイムクロッキーの集まり',
    'home.lead':
      'ロビーを作成してリンクを共有し、同期された写真とタイマーで一緒に描きましょう — 画面共有は不要です。',
    'home.createLobby': 'ロビーを作成',
    'home.createLobbyFailed': 'ロビーの作成に失敗しました',
    'home.creatingLobby': 'ロビーを作成中…',
    'home.howItWorks': 'ご利用方法',
    'home.step1.title': '作成と共有',
    'home.step1.desc':
      'ロビーを起動し、ドローインググループにリンクを送信します。',
    'home.step2.title': 'リファレンスの選択',
    'home.step2.desc': '管理者が全員で描く写真をPixabayから選択します。',
    'home.step3.title': '同期して描く',
    'home.step3.desc':
      '全員に同じ写真とカウントダウンが表示されるタイマー制ラウンドです。',

    // LobbyLayout & Page
    'lobby.badge.admin': 'ホスト',
    'lobby.badge.participant': '参加者',
    'lobby.connection.connecting': '接続中…',
    'lobby.connection.connected': '接続済み',
    'lobby.connection.reconnecting': '再接続を試みています…',
    'lobby.connection.disconnected': '接続切断',
    'lobby.connection.lost': '接続が切断されました。再接続を試みています…',
    'lobby.loadingState': 'ロビー情報を読み込み中…',
    'lobby.invalidLink': '無効なロビーリンクです。',
    'lobby.backHome': 'ホームに戻る',
    'lobby.participantCount': '{count}人の参加者',
    'lobby.participantCountPlural': '{count}人の参加者',
    'lobby.lobbyId': 'ロビー {id}',

    // CopyLobbyLinkButton
    'copy.copied': 'コピーしました！',
    'copy.failed': 'コピー失敗',
    'copy.link': 'リンクをコピー',

    // AdminControls
    'admin.startSession': 'セッション開始',
    'admin.starting': '開始中…',
    'admin.nextPhoto': '次の写真',
    'admin.loading': '読み込み中…',
    'admin.endSession': 'セッション終了',
    'admin.ending': '終了中…',
    'admin.actionFailed': '処理失敗',

    // RoundDurationPicker
    'duration.label': 'ラウンド時間',
    'duration.unit': '分',
    'duration.minutes': '1ラウンド {count} 分',
    'duration.ariaGroup': 'ラウンド時間（分）',
    'duration.decrease': 'ラウンド時間を短くする',
    'duration.increase': 'ラウンド時間を長くする',
    'duration.hint': '左右のボタンで調整できます。',
    'duration.updateFailed': 'ラウンド時間の更新に失敗しました',

    // ParticipantWaitPanel
    'wait.settingPhotos': 'リファレンス写真を設定中...',
    'wait.adminChoosing':
      '管理者が写真を選択しています。セッションが開始されるまで写真は非表示になります。',
    'wait.tipLabel': 'ヒント: {title}',
    'wait.goTip': 'ヒント {idx} に移動',
    'wait.waitingForAdmin': '管理者を待っています...',
    'wait.carouselNav': '描画のヒント',
    'wait.pauseCarousel': 'ヒントの自動切り替えを一時停止',
    'wait.resumeCarousel': 'ヒントの自動切り替えを再開',

    // PhotoReviewPanel
    'review.photosSaved': '{count}枚の写真が保存されました',
    'review.instruction':
      'サムネイルにホバーすると拡大プレビューが表示されます。選択した内容でよければ、「選択完了」を押してセッションを準備してください（写真の順序は自動的にシャッフルされ、固定されます）。',
    'review.previewAria': '保存された写真 {index} / {total} のプレビュー',
    'review.editSelection': '選択を編集',
    'review.selectionComplete': '選択完了',
    'review.shuffling': 'シャッフル中…',
    'review.modalAria': '写真プレビュー',
    'review.closePreview': 'プレビューを閉じる',
    'review.modalAlt': 'リファレンス写真 {index} / {total}',

    // PhotoSelectionPanel
    'selection.errSave': '選択内容의保存に失敗しました',
    'selection.errConfirm': '選択内容の確定に失敗しました',
    'selection.errReopen': '写真選択の再オープンに失敗しました',
    'selection.saving': '保存中…',
    'selection.saveCount': '{count}枚の写真を保存',

    // ReadyPanel
    'ready.photosReady': '枚の写真の準備完了',
    'ready.desc':
      '写真の順序はシャッフルされ、非表示になっています。各描画ラウンドが始まるまでサムネイルは表示されません。',
    'ready.hint': '管理者が開始するのを待っています…',

    // DrawingPanel
    'draw.exitFullscreen': '全画面表示の終了',
    'draw.enterFullscreen': '全画面表示',
    'draw.waitingPhoto': '写真を待っています…',
    'draw.attribution': '画像の出典:',
    'draw.round': 'ラウンド {current} / {total}',
    'draw.startsIn': 'ラウンドが {count} 秒後に始まります',
    'draw.remainingAria': '残りの描画時間',
    'draw.photoAlt': 'このクロキーラウンドのリファレンス写真',

    // SessionBreakPanels
    'break.takeBreather': 'ひと息つきましょう',
    'break.hiddenDesc':
      '次のラウンドが始まるまで、リファレンス写真は非表示になります。',
    'break.completedRound': '全 {total} ラウンド中 {current} ラウンド完了',
    'break.sessionFinished': 'セッション終了',
    'break.completedRoundsDesc':
      '全 {count} ラウンドを完了しました。皆さん、お疲れ様でした！',
    'break.completedRoundsDescPlural':
      '全 {count} ラウンドを完了しました。皆さん、お疲れ様でした！',
    'break.downloadPhotos': 'お題写真をZIPでダウンロード',

    // PixabaySearchPanel
    'search.errEmpty': '検索ワードを入力してください',
    'search.errFailed': '検索に失敗しました',
    'search.fieldLabel': 'Pixabayを検索',
    'search.placeholder': '例: ポートレート、解剖学、ジェスチャー',
    'search.sort': '並び替え',
    'search.sortPopular': '人気順',
    'search.sortLatest': '最新順',
    'search.searching': '検索中',
    'search.button': '検索',
    'search.prevPage': '前のページ',
    'search.nextPage': '次のページ',
    'search.pageIndicator': '{total}ページ中 {page}ページ',
    'search.hint': '{count}枚選択中 · {recommended}枚推奨',
    'search.ariaSelect': '写真 {id} を選択',
    'search.ariaDeselect': '写真 {id} の選択を解除',
    'search.dock.title': '選択されたリファレンス写真',
    'search.dock.remove': '写真を削除',
    'search.dock.removeAria': '写真 {id} を削除',
  },
  pl: {
    // HomePage
    'home.eyebrow': 'Spotkania kroquis w czasie rzeczywistym',
    'home.lead':
      'Utwórz pokój, udostępnij link i rysujcie wspólnie ze zsynchronizowanymi zdjęciami i licznikami — udostępnianie ekranu nie jest wymagane.',
    'home.createLobby': 'Utwórz pokój',
    'home.createLobbyFailed': 'Nie udało się utworzyć pokoju',
    'home.creatingLobby': 'Tworzenie pokoju…',
    'home.howItWorks': 'Jak to działa',
    'home.step1.title': 'Utwórz i udostępnij',
    'home.step1.desc':
      'Uruchom pokój i wyślij link do swojej grupy rysunkowej.',
    'home.step2.title': 'Wybierz referencje',
    'home.step2.desc':
      'Administrator wybiera zdjęcia z Pixabay, które wszyscy będą rysować.',
    'home.step3.title': 'Rysuj w synchronizacji',
    'home.step3.desc':
      'Rundy na czas z tym samym zdjęciem i odliczaniem dla wszystkich.',

    // LobbyLayout & Page
    'lobby.badge.admin': 'Admin',
    'lobby.badge.participant': 'Uczestnik',
    'lobby.connection.connecting': 'Łączenie…',
    'lobby.connection.connected': 'Na żywo',
    'lobby.connection.reconnecting': 'Ponowne łączenie…',
    'lobby.connection.disconnected': 'Rozłączono',
    'lobby.connection.lost': 'Połączenie utracone. Próba ponownego połączenia…',
    'lobby.loadingState': 'Ładowanie stanu pokoju…',
    'lobby.invalidLink': 'Nieprawidłowy link do pokoju.',
    'lobby.backHome': 'Powrót do strony głównej',
    'lobby.participantCount': '{count} uczestnik',
    'lobby.participantCountPlural': '{count} uczestników',
    'lobby.lobbyId': 'Pokój {id}',

    // CopyLobbyLinkButton
    'copy.copied': 'Skopiowano!',
    'copy.failed': 'Błąd kopiowania',
    'copy.link': 'Kopiuj link',

    // AdminControls
    'admin.startSession': 'Rozpocznij sesję',
    'admin.starting': 'Rozpoczynanie…',
    'admin.nextPhoto': 'Następne zdjęcie',
    'admin.loading': 'Ładowanie…',
    'admin.endSession': 'Zakończ sesję',
    'admin.ending': 'Kończenie…',
    'admin.actionFailed': 'Akcja nie powiodła się',

    // RoundDurationPicker
    'duration.label': 'Czas rundy',
    'duration.unit': 'min',
    'duration.minutes': '{count} min na rundę',
    'duration.ariaGroup': 'Czas rundy w minutach',
    'duration.decrease': 'Skróć czas rundy',
    'duration.increase': 'Wydłuż czas rundy',
    'duration.hint': 'Użyj lewego i prawego przycisku, aby dostosować.',
    'duration.updateFailed': 'Nie udało się zaktualizować czasu rundy',

    // ParticipantWaitPanel
    'wait.settingPhotos': 'Konfigurowanie zdjęć referencyjnych...',
    'wait.adminChoosing':
      'Administrator wybiera zdjęcia. Pozostaną one ukryte do momentu rozpoczęcia sesji.',
    'wait.tipLabel': 'Wskazówka: {title}',
    'wait.goTip': 'Przejdź do wskazówki {idx}',
    'wait.waitingForAdmin': 'Oczekiwanie na administratora...',
    'wait.carouselNav': 'Wskazówki do rysowania',
    'wait.pauseCarousel': 'Wstrzymaj rotację wskazówek',
    'wait.resumeCarousel': 'Wznów rotację wskazówek',

    // PhotoReviewPanel
    'review.photosSaved': 'Zapisano {count} zdjęć',
    'review.instruction':
      'Najedź kursorem na miniaturkę, aby podglądnąć ją w pełnym rozmiarze. Jeśli zestaw Ci odpowiada, kliknij „Wybór zakończony”, aby przygotować sesję (kolejność zdjęć zostanie automatycznie wymieszana).',
    'review.previewAria': 'Podgląd zapisanego zdjęcia {index} z {total}',
    'review.editSelection': 'Edytuj wybór',
    'review.selectionComplete': 'Wybór zakończony',
    'review.shuffling': 'Mieszanie…',
    'review.modalAria': 'Podgląd zdjęcia',
    'review.closePreview': 'Zamknij podgląd',
    'review.modalAlt': 'Zdjęcie referencyjne {index} z {total}',

    // PhotoSelectionPanel
    'selection.errSave': 'Nie udało się zapisać wyboru',
    'selection.errConfirm': 'Nie udało się potwierdzić wyboru',
    'selection.errReopen': 'Nie udało się ponownie otworzyć wyboru zdjęć',
    'selection.saving': 'Zapisywanie…',
    'selection.saveCount': 'Zapisz {count} zdjęć',

    // ReadyPanel
    'ready.photosReady': 'zdjęć gotowych',
    'ready.desc':
      'Kolejność jest wymieszana i ukryta. Miniaturki pozostają wyłączone do czasu rozpoczęcia każdej rundy rysowania.',
    'ready.hint': 'Oczekiwanie na rozpoczęcie przez administratora…',

    // DrawingPanel
    'draw.exitFullscreen': 'Wyjdź z pełnego ekranu',
    'draw.enterFullscreen': 'Pełny ekran',
    'draw.waitingPhoto': 'Oczekiwanie na zdjęcie…',
    'draw.attribution': 'Obraz z',
    'draw.round': 'Runda {current} / {total}',
    'draw.startsIn': 'Runda rozpoczyna się za {count} sekund',
    'draw.remainingAria': 'Pozostały czas rysowania',
    'draw.photoAlt': 'Zdjęcie referencyjne na tę rundę kroquis',

    // SessionBreakPanels
    'break.takeBreather': 'Złap oddech',
    'break.hiddenDesc':
      'Zdjęcie referencyjne jest ukryte do momentu rozpoczęcia następnej rundy.',
    'break.completedRound': 'Ukończono rundę {current} z {total}',
    'break.sessionFinished': 'Sesja zakończona',
    'break.completedRoundsDesc':
      'Ukończono {count} rundę. Dobra robota, wszyscy.',
    'break.completedRoundsDescPlural':
      'Ukończono {count} rund. Dobra robota, wszyscy.',
    'break.downloadPhotos': 'Pobierz zdjęcia referencyjne (ZIP)',

    // PixabaySearchPanel
    'search.errEmpty': 'Wpisz wyszukiwane hasło',
    'search.errFailed': 'Wyszukiwanie nie powiodło się',
    'search.fieldLabel': 'Szukaj w Pixabay',
    'search.placeholder': 'np. portret, anatomia, gest',
    'search.sort': 'Sortuj',
    'search.sortPopular': 'Popularne',
    'search.sortLatest': 'Najnowsze',
    'search.searching': 'Wyszukiwanie',
    'search.button': 'Szukaj',
    'search.prevPage': 'Poprzednia strona',
    'search.nextPage': 'Następna strona',
    'search.pageIndicator': 'Strona {page} z {total}',
    'search.hint': 'Wybrano {count} · Zalecane {recommended}',
    'search.ariaSelect': 'Wybierz zdjęcie {id}',
    'search.ariaDeselect': 'Odznacz zdjęcie {id}',
    'search.dock.title': 'Wybrane zdjęcia referencyjne',
    'search.dock.remove': 'Usuń zdjęcie',
    'search.dock.removeAria': 'Usuń zdjęcie {id}',
  },
  zh: {
    // HomePage
    'home.eyebrow': '实时速写聚会',
    'home.lead':
      '创建大厅、分享链接，即可同步照片和计时器一起画画——无需屏幕共享。',
    'home.createLobby': '创建大厅',
    'home.createLobbyFailed': '创建大厅失败',
    'home.creatingLobby': '正在创建大厅…',
    'home.howItWorks': '使用方法',
    'home.step1.title': '创建并分享',
    'home.step1.desc': '创建大厅，把链接发给绘画小组。',
    'home.step2.title': '选择参考图',
    'home.step2.desc': '管理员从 Pixabay 挑选大家一起画的图片。',
    'home.step3.title': '同步作画',
    'home.step3.desc': '定时回合，所有人看到相同的图片和倒计时。',

    // LobbyLayout & Page
    'lobby.badge.admin': '管理员',
    'lobby.badge.participant': '参与者',
    'lobby.connection.connecting': '连接中…',
    'lobby.connection.connected': '已连接',
    'lobby.connection.reconnecting': '重新连接中…',
    'lobby.connection.disconnected': '已断开',
    'lobby.connection.lost': '连接已断开，正在尝试重新连接…',
    'lobby.loadingState': '正在加载大厅状态…',
    'lobby.invalidLink': '无效的大厅链接。',
    'lobby.backHome': '返回首页',
    'lobby.participantCount': '{count} 位参与者',
    'lobby.participantCountPlural': '{count} 位参与者',
    'lobby.lobbyId': '大厅 {id}',

    // CopyLobbyLinkButton
    'copy.copied': '已复制！',
    'copy.failed': '复制失败',
    'copy.link': '复制链接',

    // AdminControls
    'admin.startSession': '开始会话',
    'admin.starting': '开始中…',
    'admin.nextPhoto': '下一张',
    'admin.loading': '加载中…',
    'admin.endSession': '结束会话',
    'admin.ending': '结束中…',
    'admin.actionFailed': '操作失败',

    // RoundDurationPicker
    'duration.label': '回合时长',
    'duration.unit': '分钟',
    'duration.minutes': '每回合 {count} 分钟',
    'duration.ariaGroup': '回合时长（分钟）',
    'duration.decrease': '缩短回合时长',
    'duration.increase': '延长回合时长',
    'duration.hint': '使用左右按钮调节。',
    'duration.updateFailed': '更新回合时长失败',

    // ParticipantWaitPanel
    'wait.settingPhotos': '正在准备参考图片…',
    'wait.adminChoosing': '管理员正在挑选图片，会话开始前不会显示。',
    'wait.tipLabel': '提示：{title}',
    'wait.goTip': '跳转到提示 {idx}',
    'wait.waitingForAdmin': '等待管理员…',
    'wait.carouselNav': '绘画提示',
    'wait.pauseCarousel': '暂停提示轮播',
    'wait.resumeCarousel': '继续提示轮播',

    // PhotoReviewPanel
    'review.photosSaved': '已保存 {count} 张图片',
    'review.instruction':
      '鼠标悬停缩略图可预览大图。确认无误后点击完成选择以准备会话（图片顺序会自动打乱）。',
    'review.previewAria': '预览已保存的图片 {index} / {total}',
    'review.editSelection': '修改选择',
    'review.selectionComplete': '完成选择',
    'review.shuffling': '打乱中…',
    'review.modalAria': '图片预览',
    'review.closePreview': '关闭预览',
    'review.modalAlt': '参考图片 {index} / {total}',

    // PhotoSelectionPanel
    'selection.errSave': '保存选择失败',
    'selection.errConfirm': '确认选择失败',
    'selection.errReopen': '重新打开图片选择失败',
    'selection.saving': '保存中…',
    'selection.saveCount': '保存 {count} 张图片',

    // ReadyPanel
    'ready.photosReady': '张图片已就绪',
    'ready.desc': '图片顺序已打乱并隐藏，每个作画回合开始前不会显示缩略图。',
    'ready.hint': '等待管理员开始…',

    // DrawingPanel
    'draw.exitFullscreen': '退出全屏',
    'draw.enterFullscreen': '全屏',
    'draw.waitingPhoto': '等待图片…',
    'draw.attribution': '图片来源：',
    'draw.round': '第 {current} / {total} 回合',
    'draw.startsIn': '{count} 秒后开始回合',
    'draw.remainingAria': '剩余作画时间',
    'draw.photoAlt': '本回合速写参考图片',

    // SessionBreakPanels
    'break.takeBreather': '休息一下',
    'break.hiddenDesc': '参考图片已隐藏，直到下一回合开始。',
    'break.completedRound': '已完成第 {current} / {total} 回合',
    'break.sessionFinished': '会话结束',
    'break.completedRoundsDesc': '完成了 {count} 个回合，大家辛苦了！',
    'break.completedRoundsDescPlural': '完成了 {count} 个回合，大家辛苦了！',
    'break.downloadPhotos': '下载参考图片 (ZIP)',

    // PixabaySearchPanel
    'search.errEmpty': '请输入搜索词',
    'search.errFailed': '搜索失败',
    'search.fieldLabel': '搜索 Pixabay',
    'search.placeholder': '例如：portrait, anatomy, gesture',
    'search.sort': '排序',
    'search.sortPopular': '热门',
    'search.sortLatest': '最新',
    'search.searching': '搜索中',
    'search.button': '搜索',
    'search.prevPage': '上一页',
    'search.nextPage': '下一页',
    'search.pageIndicator': '第 {page} / {total} 页',
    'search.hint': '已选 {count} 张 · 建议 {recommended} 张',
    'search.ariaSelect': '选择图片 {id}',
    'search.ariaDeselect': '取消选择图片 {id}',
    'search.dock.title': '已选参考图片',
    'search.dock.remove': '移除图片',
    'search.dock.removeAria': '移除图片 {id}',
  },
}

export function t(
  key: string,
  params?: Record<string, string | number>,
): string {
  const lang = getCurrentLanguage()
  const template = TRANSLATIONS[lang]?.[key] || TRANSLATIONS.en[key] || key
  if (!params) return template
  return Object.entries(params).reduce(
    (acc, [k, v]) => acc.replace(new RegExp(`{${k}}`, 'g'), String(v)),
    template,
  )
}

// Structured Localized Phase Messages
export interface PhaseMessage {
  title: string
  description: string
}

const LOCALIZED_PHASES: Record<Language, Record<LobbyPhase, PhaseMessage>> = {
  en: {
    WAITING: {
      title: 'Waiting for the host',
      description: 'The admin is getting things ready. Stay on this page.',
    },
    SELECTING: {
      title: 'Reviewing reference photos',
      description:
        'The admin is checking the saved photos before shuffling and starting the session.',
    },
    READY: {
      title: 'Ready to start',
      description:
        'Photos are shuffled and hidden. Waiting for the admin to begin.',
    },
    DRAWING: {
      title: 'Draw time',
      description:
        'Focus on the reference photo. The timer is server-controlled.',
    },
    BETWEEN_ROUNDS: {
      title: 'Round break',
      description: 'Take a short break before the next pose.',
    },
    FINISHED: {
      title: 'Session complete',
      description: 'Thanks for drawing together. See you next week.',
    },
  },
  ko: {
    WAITING: {
      title: '방장 기다리는 중',
      description: '방장이 준비 중이에요. 이 페이지에서 기다려 주세요.',
    },
    SELECTING: {
      title: '레퍼런스 사진 확인 중',
      description: '방장이 저장된 사진을 확인한 뒤 섞고 세션을 시작합니다.',
    },
    READY: {
      title: '시작 준비 완료',
      description:
        '사진 순서를 섞어 숨겨 두었어요. 방장이 시작할 때까지 기다려 주세요.',
    },
    DRAWING: {
      title: '그리기 시간!',
      description: '레퍼런스 사진에 집중하세요. 타이머는 서버가 맞춰 줍니다.',
    },
    BETWEEN_ROUNDS: {
      title: '라운드 쉬는 시간',
      description: '다음 포즈 전에 잠깐 쉬어 가세요.',
    },
    FINISHED: {
      title: '세션 종료',
      description: '함께 그려 주셔서 고마워요. 다음 주에 또 만나요!',
    },
  },
  ja: {
    WAITING: {
      title: 'ホストの待機中',
      description:
        '管理者が準備を進めています。このページのままお待ちください。',
    },
    SELECTING: {
      title: 'リファレンス写真の確認中',
      description:
        '管理者が写真をシャッフルしてセッションを開始する前に、保存された写真を確認しています。',
    },
    READY: {
      title: '開始準備完了',
      description:
        '写真はシャッフルされ、非表示になっています。管理者が開始するのをお待ちください。',
    },
    DRAWING: {
      title: '描画時間',
      description:
        'リファレンス写真に集中してください。タイマーはサーバーによって制御されています。',
    },
    BETWEEN_ROUNDS: {
      title: 'ラウンド間の休憩',
      description: '次のポーズが始まる前に、少し休憩してください。',
    },
    FINISHED: {
      title: 'セッション終了',
      description:
        '一緒に描いていただきありがとうございました。また来週お会いしましょう。',
    },
  },
  pl: {
    WAITING: {
      title: 'Oczekiwanie na hosta',
      description: 'Administrator przygotowuje pokój. Pozostań na tej stronie.',
    },
    SELECTING: {
      title: 'Przeglądanie zdjęć referencyjnych',
      description:
        'Administrator sprawdza zapisane zdjęcia przed ich wymieszaniem i rozpoczęciem sesji.',
    },
    READY: {
      title: 'Gotowy do rozpoczęcia',
      description:
        'Zdjęcia są wymieszane i ukryte. Oczekiwanie na rozpoczęcie przez administratora.',
    },
    DRAWING: {
      title: 'Czas na rysowanie',
      description:
        'Skoncentruj się na zdjęciu referencyjnym. Licznik czasu jest kontrolowany przez serwer.',
    },
    BETWEEN_ROUNDS: {
      title: 'Przerwa między rundami',
      description: 'Zrób sobie krótką przerwę przed następną pozą.',
    },
    FINISHED: {
      title: 'Sesja zakończona',
      description:
        'Dziękujemy za wspólne rysowanie. Do zobaczenia w przyszłym tygodniu.',
    },
  },
  zh: {
    WAITING: {
      title: '等待管理员',
      description: '管理员正在准备，请留在此页面。',
    },
    SELECTING: {
      title: '确认参考图片',
      description: '管理员正在检查已保存的图片，随后会打乱顺序并开始会话。',
    },
    READY: {
      title: '准备就绪',
      description: '图片已打乱并隐藏，等待管理员开始。',
    },
    DRAWING: {
      title: '作画时间',
      description: '专注参考图片，计时器由服务器控制。',
    },
    BETWEEN_ROUNDS: {
      title: '回合间歇',
      description: '下一个姿势开始前，先放松一下。',
    },
    FINISHED: {
      title: '会话结束',
      description: '感谢一起作画，下周再见！',
    },
  },
}

export function getPhaseMessage(phase: LobbyPhase): PhaseMessage {
  const lang = getCurrentLanguage()
  return LOCALIZED_PHASES[lang]?.[phase] || LOCALIZED_PHASES.en[phase]
}

export function getConnectionLabel(
  status: 'connecting' | 'connected' | 'reconnecting' | 'disconnected',
): string {
  return t(`lobby.connection.${status}`) || status
}

// Localized Tips
export interface Tip {
  title: string
  text: string
}

const LOCALIZED_TIPS: Record<Language, Tip[]> = {
  en: [
    {
      title: 'Gesture First',
      text: 'Capture the main line of action and body flow first, ignoring all small details.',
    },
    {
      title: 'Squint to Simplify',
      text: 'Squinting helps block out high-frequency details so you can see the main shapes of light and shadow.',
    },
    {
      title: 'Draw with your Arm',
      text: 'Use your elbow and shoulder to make long, sweeping strokes instead of scratchy wrist movements.',
    },
    {
      title: 'Exaggerate the Pose',
      text: 'In short gesture sketches, it is better to draw a pose that looks more dynamic and expressive than the photo.',
    },
    {
      title: 'Negative Spaces',
      text: 'Look at the empty shapes formed between limbs to gauge proportions and angles accurately.',
    },
    {
      title: 'Proportions over Details',
      text: 'Leave hands, feet, and faces blank until the main torso and overall skeleton structure are locked in.',
    },
    {
      title: 'Keep it Loose',
      text: 'Remember, this is a croquis! Speed and feeling are more important than a polished, perfect drawing.',
    },
  ],
  ko: [
    {
      title: '제스처 먼저',
      text: '세부 묘사는 나중에. 몸의 흐름과 동세부터 잡으세요.',
    },
    {
      title: '눈 가늘게 뜨고 보기',
      text: '눈을 살짝 감으면 자잘한 디테일이 묻혀서, 명암의 큰 덩어리가 더 잘 보입니다.',
    },
    {
      title: '팔 전체로 그리기',
      text: '손목만 쓰지 말고 팔꿈치·어깨까지 써서 길고 시원한 선을 그리세요.',
    },
    {
      title: '포즈 과장하기',
      text: '짧은 크로키에서는 사진보다 포즈를 더 과장해서 그리는 편이 낫습니다.',
    },
    {
      title: '여백(네거티브 스페이스) 보기',
      text: '팔다리 사이 빈 공간 모양을 보면 비율과 각도를 더 정확히 잴 수 있어요.',
    },
    {
      title: '디테일보다 비율',
      text: '몸통과 척추가 잡힐 때까지 손·발·얼굴 같은 디테일은 비워 두세요.',
    },
    {
      title: '느슨하게 그리기',
      text: '크로키예요! 완벽하게 다듬기보다 속도와 느낌이 더 중요합니다.',
    },
  ],
  ja: [
    {
      title: 'ジェスチャー優先',
      text: '細かいディテールは無視し、体の流れや主要な動勢を最初に捉えましょう。',
    },
    {
      title: '目を細めて単純化する',
      text: '目を少し細めることで細かいディテールが排除され、光と影の大きな塊を捉えやすくなります。',
    },
    {
      title: '腕全体で描く',
      text: '手首だけで細かく描くのではなく、肘や肩を使って長く伸びやかな線を描きましょう。',
    },
    {
      title: 'ポーズを誇張する',
      text: '短い時間でのクロッキーでは、実際の写真よりもポーズをよりダイナミックで表現豊かに描く方が効果的です。',
    },
    {
      title: 'ネガティブスペース（余白）の活用',
      text: '手足の間にできる余白の形を観察すると、比率や角度をより正確に把握できます。',
    },
    {
      title: 'ディテールより比率優先',
      text: '胴体と脊椎の構造がしっかりと描けるまでは、手、足、顔などの細かい描写は控えましょう。',
    },
    {
      title: 'リラックスして描く',
      text: 'これはクロッキーです！完璧に仕上げることよりも、スピード感とフィーリングを捉えることがはるかに重要です。',
    },
  ],
  pl: [
    {
      title: 'Gest przede wszystkim',
      text: 'Najpierw uchwyć główną linię ruchu i przepływ ciała, ignorując drobne szczegóły.',
    },
    {
      title: 'Zmruż oczy, aby uprościć',
      text: 'Mrużenie oczu pomaga odfiltrować drobne szczegóły, pozwalając dostrzec główne kształty światła i cienia.',
    },
    {
      title: 'Rysuj całą ręką',
      text: 'Używaj łokcia i ramienia do wykonywania długich, zamaszystych pociągnięć zamiast drobnych ruchów nadgarstka.',
    },
    {
      title: 'Przesadzaj z pozą',
      text: 'W krótkich szkicach gestów lepiej narysować pozę, która wygląda na bardziej dynamiczną i ekspresyjną niż na zdjęciu.',
    },
    {
      title: 'Negatywne przestrzenie',
      text: 'Przyjrzyj się pustym kształtom tworzącym się między kończynami, aby dokładnie ocenić proporcje i kąty.',
    },
    {
      title: 'Proporcje ponad szczegółami',
      text: 'Pozostaw dłonie, stopy i twarze puste, dopóki główny tułów i ogólna struktura szkieletu nie zostaną ustalone.',
    },
    {
      title: 'Rysuj swobodnie',
      text: 'Pamiętaj, to jest kroquis! Szybkość i wyczucie są ważniejsze niż dopracowany, idealny rysunek.',
    },
  ],
  zh: [
    {
      title: '先抓动态',
      text: '先捕捉动作主线和身体的流动感，暂时忽略细节。',
    },
    {
      title: '眯眼简化',
      text: '眯起眼睛可以过滤细碎细节，更容易看到明暗的大块形状。',
    },
    {
      title: '用整条手臂画',
      text: '用肘部和肩部画长而流畅的线，而不是只用腕部来回蹭。',
    },
    {
      title: '夸张姿势',
      text: '短时间速写时，把姿势画得比照片更有动感和表现力往往更好。',
    },
    {
      title: '负形空间',
      text: '观察四肢之间的空白形状，能更准确判断比例和角度。',
    },
    {
      title: '比例优先于细节',
      text: '躯干和骨架结构确定之前，手、脚、脸等细节可以先留白。',
    },
    {
      title: '放松地画',
      text: '这是速写！速度和感觉比精致完美的成品更重要。',
    },
  ],
}

export function getLocalizedTips(): Tip[] {
  const lang = getCurrentLanguage()
  return LOCALIZED_TIPS[lang] || LOCALIZED_TIPS.en
}
