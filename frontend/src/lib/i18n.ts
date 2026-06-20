import type { LobbyPhase } from '../types/lobby'

export type Language = 'ko' | 'ja' | 'pl' | 'en'

export function getBrowserLanguage(): Language {
  const lang = navigator.language.toLowerCase()
  if (lang.startsWith('ko')) return 'ko'
  if (lang.startsWith('ja')) return 'ja'
  if (lang.startsWith('pl')) return 'pl'
  return 'en'
}

// Simple translation lookup with placeholder support like {count}
const TRANSLATIONS: Record<Language, Record<string, string>> = {
  en: {
    // HomePage
    'home.eyebrow': 'Real-time croquis meetups',
    'home.lead': 'Create a lobby, share the link, and draw together with synchronized photos and timers — no screen sharing required.',
    'home.createLobby': 'Create lobby',
    'home.createLobbyFailed': 'Failed to create lobby',
    'home.creatingLobby': 'Creating lobby…',
    'home.howItWorks': 'How it works',
    'home.step1.title': 'Create & share',
    'home.step1.desc': 'Start a lobby and send the link to your drawing group.',
    'home.step2.title': 'Pick references',
    'home.step2.desc': 'The admin selects photos from Pixabay for everyone to draw.',
    'home.step3.title': 'Draw in sync',
    'home.step3.desc': 'Timed rounds with the same photo and countdown for all.',
    
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
    
    // ParticipantWaitPanel
    'wait.settingPhotos': 'Setting up reference photos...',
    'wait.adminChoosing': 'The admin is choosing photos. They stay hidden until the session starts.',
    'wait.tipLabel': 'Tip: {title}',
    'wait.goTip': 'Go to tip {idx}',
    'wait.waitingForAdmin': 'Waiting for admin...',
    
    // PhotoReviewPanel
    'review.photosSaved': '{count} photos saved',
    'review.instruction': 'Hover a thumbnail to preview it at full size. When you are happy with the set, shuffle and lock the order to start.',
    'review.previewAria': 'Preview saved photo {index} of {total}',
    'review.editSelection': 'Edit selection',
    'review.selectionComplete': 'Selection complete',
    'review.shuffling': 'Shuffling…',
    
    // PhotoSelectionPanel
    'selection.errSave': 'Failed to save selection',
    'selection.errConfirm': 'Failed to confirm selection',
    'selection.errReopen': 'Failed to reopen photo selection',
    'selection.saving': 'Saving…',
    'selection.saveCount': 'Save {count} photos',
    
    // ReadyPanel
    'ready.photosReady': 'photos ready',
    'ready.desc': 'The order is shuffled and hidden. Thumbnails stay off until each draw round begins.',
    'ready.hint': 'Waiting for the admin to start…',
    
    // DrawingPanel
    'draw.exitFullscreen': 'Exit Fullscreen',
    'draw.enterFullscreen': 'Enter Fullscreen',
    'draw.waitingPhoto': 'Waiting for photo…',
    'draw.attribution': 'Image from',
    'draw.round': 'Round {current} / {total}',
    'draw.startsIn': 'Round starts in {count} seconds',
    'draw.remainingAria': 'Draw time remaining',
    
    // SessionBreakPanels
    'break.takeBreather': 'Take a breather',
    'break.hiddenDesc': 'The reference photo is hidden until the next round starts.',
    'break.completedRound': 'Completed round {current} of {total}',
    'break.sessionFinished': 'Session finished',
    'break.completedRoundsDesc': 'You completed {count} round. Great work everyone.',
    'break.completedRoundsDescPlural': 'You completed {count} rounds. Great work everyone.',
    
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
  },
  ko: {
    // HomePage
    'home.eyebrow': '실시간 크로키 모임',
    'home.lead': '로비를 만들고 링크를 공유하여 동기화된 사진과 타이머로 함께 그리세요 — 화면 공유가 필요 없습니다.',
    'home.createLobby': '로비 만들기',
    'home.createLobbyFailed': '로비 생성 실패',
    'home.creatingLobby': '로비 만드는 중…',
    'home.howItWorks': '이용 방법',
    'home.step1.title': '방 생성 및 공유',
    'home.step1.desc': '로비를 생성하고 드로잉 그룹에 링크를 보내세요.',
    'home.step2.title': '레퍼런스 선택',
    'home.step2.desc': '방장이 Pixabay에서 모두가 그릴 사진을 선택합니다.',
    'home.step3.title': '동시에 그리기',
    'home.step3.desc': '모두에게 동일한 사진과 카운트다운이 제공되는 타이머 라운드입니다.',
    
    // LobbyLayout & Page
    'lobby.badge.admin': '방장',
    'lobby.badge.participant': '참가자',
    'lobby.connection.connecting': '연결 중…',
    'lobby.connection.connected': '연결됨',
    'lobby.connection.reconnecting': '재연결 중…',
    'lobby.connection.disconnected': '연결 끊김',
    'lobby.connection.lost': '연결이 끊겼습니다. 재연결을 시도 중…',
    'lobby.loadingState': '로비 상태 로딩 중…',
    'lobby.invalidLink': '올바르지 않은 로비 링크입니다.',
    'lobby.backHome': '홈으로 돌아가기',
    'lobby.participantCount': '{count}명 참가 중',
    'lobby.participantCountPlural': '{count}명 참가 중',
    'lobby.lobbyId': '로비 {id}',
    
    // CopyLobbyLinkButton
    'copy.copied': '복사 완료!',
    'copy.failed': '복사 실패',
    'copy.link': '링크 복사',
    
    // AdminControls
    'admin.startSession': '세션 시작',
    'admin.starting': '시작 중…',
    'admin.nextPhoto': '다음 사진',
    'admin.loading': '로딩 중…',
    'admin.endSession': '세션 종료',
    'admin.ending': '종료 중…',
    'admin.actionFailed': '작업 실패',
    
    // ParticipantWaitPanel
    'wait.settingPhotos': '레퍼런스 사진 설정 중...',
    'wait.adminChoosing': '방장이 사진을 선택하는 중입니다. 세션이 시작될 때까지 사진은 숨겨집니다.',
    'wait.tipLabel': '팁: {title}',
    'wait.goTip': '팁 {idx}으로 이동',
    'wait.waitingForAdmin': '방장을 기다리는 중...',
    
    // PhotoReviewPanel
    'review.photosSaved': '{count}장의 사진이 저장됨',
    'review.instruction': '썸네일에 마우스를 올리면 크게 미리 볼 수 있습니다. 사진 세트가 마음에 들면 순서를 섞고 고정한 뒤 시작해 주세요.',
    'review.previewAria': '저장된 사진 {index} / {total} 미리보기',
    'review.editSelection': '선택 수정',
    'review.selectionComplete': '선택 완료',
    'review.shuffling': '셔플 중…',
    
    // PhotoSelectionPanel
    'selection.errSave': '선택 사항 저장 실패',
    'selection.errConfirm': '선택 사항 확정 실패',
    'selection.errReopen': '사진 선택 재오픈 실패',
    'selection.saving': '저장 중…',
    'selection.saveCount': '{count}장의 사진 저장',
    
    // ReadyPanel
    'ready.photosReady': '장의 사진 준비됨',
    'ready.desc': '사진 순서가 섞인 뒤 숨겨졌습니다. 각 드로잉 라운드가 시작되기 전까지 썸네일은 보이지 않습니다.',
    'ready.hint': '방장이 시작하기를 기다리는 중…',
    
    // DrawingPanel
    'draw.exitFullscreen': '전체화면 종료',
    'draw.enterFullscreen': '전체화면',
    'draw.waitingPhoto': '사진을 기다리는 중…',
    'draw.attribution': '사진 출처:',
    'draw.round': '라운드 {current} / {total}',
    'draw.startsIn': '라운드가 {count}초 후 시작됩니다',
    'draw.remainingAria': '남은 그리기 시간',
    
    // SessionBreakPanels
    'break.takeBreather': '잠시 휴식을 취하세요',
    'break.hiddenDesc': '다음 라운드가 시작될 때까지 레퍼런스 사진은 숨겨집니다.',
    'break.completedRound': '{total}라운드 중 {current}라운드 완료',
    'break.sessionFinished': '세션 종료',
    'break.completedRoundsDesc': '총 {count}라운드를 완료했습니다. 모두 수고하셨습니다!',
    'break.completedRoundsDescPlural': '총 {count}라운드를 완료했습니다. 모두 수고하셨습니다!',
    
    // PixabaySearchPanel
    'search.errEmpty': '검색어를 입력하세요',
    'search.errFailed': '검색 실패',
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
    'search.hint': '{count}장 선택됨 · {recommended}장 권장',
    'search.ariaSelect': '사진 {id} 선택',
    'search.ariaDeselect': '사진 {id} 선택 해제',
    'search.dock.title': '선택된 레퍼런스 사진',
    'search.dock.remove': '사진 제거',
  },
  ja: {
    // HomePage
    'home.eyebrow': 'リアルタイムクロッキーの集まり',
    'home.lead': 'ロビーを作成してリンクを共有し、同期された写真とタイマーで一緒に描きましょう — 画面共有は不要です。',
    'home.createLobby': 'ロビーを作成',
    'home.createLobbyFailed': 'ロビーの作成に失敗しました',
    'home.creatingLobby': 'ロビーを作成中…',
    'home.howItWorks': 'ご利用方法',
    'home.step1.title': '作成と共有',
    'home.step1.desc': 'ロビーを起動し、ドローインググループにリンクを送信します。',
    'home.step2.title': 'リファレンスの選択',
    'home.step2.desc': '管理者が全員で描く写真をPixabayから選択します。',
    'home.step3.title': '同期して描く',
    'home.step3.desc': '全員に同じ写真とカウントダウンが表示されるタイマー制ラウンドです。',
    
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
    
    // ParticipantWaitPanel
    'wait.settingPhotos': 'リファレンス写真を設定中...',
    'wait.adminChoosing': '管理者が写真を選択しています。セッションが開始されるまで写真は非表示になります。',
    'wait.tipLabel': 'ヒント: {title}',
    'wait.goTip': 'ヒント {idx} に移動',
    'wait.waitingForAdmin': '管理者を待っています...',
    
    // PhotoReviewPanel
    'review.photosSaved': '{count}枚の写真が保存されました',
    'review.instruction': 'サムネイルにホバーすると拡大プレビューが表示されます。選択した内容でよければ、シャッフルして順序を固定し、開始してください。',
    'review.previewAria': '保存された写真 {index} / {total} のプレビュー',
    'review.editSelection': '選択を編集',
    'review.selectionComplete': '選択完了',
    'review.shuffling': 'シャッフル中…',
    
    // PhotoSelectionPanel
    'selection.errSave': '選択内容의保存に失敗しました',
    'selection.errConfirm': '選択内容の確定に失敗しました',
    'selection.errReopen': '写真選択の再オープンに失敗しました',
    'selection.saving': '保存中…',
    'selection.saveCount': '{count}枚の写真を保存',
    
    // ReadyPanel
    'ready.photosReady': '枚の写真の準備完了',
    'ready.desc': '写真の順序はシャッフルされ、非表示になっています。各描画ラウンドが始まるまでサムネイルは表示されません。',
    'ready.hint': '管理者が開始するのを待っています…',
    
    // DrawingPanel
    'draw.exitFullscreen': '全画面表示の終了',
    'draw.enterFullscreen': '全画面表示',
    'draw.waitingPhoto': '写真を待っています…',
    'draw.attribution': '画像の出典:',
    'draw.round': 'ラウンド {current} / {total}',
    'draw.startsIn': 'ラウンドが {count} 秒後に始まります',
    'draw.remainingAria': '残りの描画時間',
    
    // SessionBreakPanels
    'break.takeBreather': 'ひと息つきましょう',
    'break.hiddenDesc': '次のラウンドが始まるまで、リファレンス写真は非表示になります。',
    'break.completedRound': '全 {total} ラウンド中 {current} ラウンド完了',
    'break.sessionFinished': 'セッション終了',
    'break.completedRoundsDesc': '全 {count} ラウンドを完了しました。皆さん、お疲れ様でした！',
    'break.completedRoundsDescPlural': '全 {count} ラウンドを完了しました。皆さん、お疲れ様でした！',
    
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
  },
  pl: {
    // HomePage
    'home.eyebrow': 'Spotkania kroquis w czasie rzeczywistym',
    'home.lead': 'Utwórz pokój, udostępnij link i rysujcie wspólnie ze zsynchronizowanymi zdjęciami i licznikami — udostępnianie ekranu nie jest wymagane.',
    'home.createLobby': 'Utwórz pokój',
    'home.createLobbyFailed': 'Nie udało się utworzyć pokoju',
    'home.creatingLobby': 'Tworzenie pokoju…',
    'home.howItWorks': 'Jak to działa',
    'home.step1.title': 'Utwórz i udostępnij',
    'home.step1.desc': 'Uruchom pokój i wyślij link do swojej grupy rysunkowej.',
    'home.step2.title': 'Wybierz referencje',
    'home.step2.desc': 'Administrator wybiera zdjęcia z Pixabay, które wszyscy będą rysować.',
    'home.step3.title': 'Rysuj w synchronizacji',
    'home.step3.desc': 'Rundy na czas z tym samym zdjęciem i odliczaniem dla wszystkich.',
    
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
    
    // ParticipantWaitPanel
    'wait.settingPhotos': 'Konfigurowanie zdjęć referencyjnych...',
    'wait.adminChoosing': 'Administrator wybiera zdjęcia. Pozostaną one ukryte do momentu rozpoczęcia sesji.',
    'wait.tipLabel': 'Wskazówka: {title}',
    'wait.goTip': 'Przejdź do wskazówki {idx}',
    'wait.waitingForAdmin': 'Oczekiwanie na administratora...',
    
    // PhotoReviewPanel
    'review.photosSaved': 'Zapisano {count} zdjęć',
    'review.instruction': 'Najedź kursorem na miniaturkę, aby podglądnąć ją w pełnym rozmiarze. Kiedy zestaw Ci się spodoba, wymieszaj i zablokuj kolejność, aby rozpocząć.',
    'review.previewAria': 'Podgląd zapisanego zdjęcia {index} z {total}',
    'review.editSelection': 'Edytuj wybór',
    'review.selectionComplete': 'Wybór zakończony',
    'review.shuffling': 'Mieszanie…',
    
    // PhotoSelectionPanel
    'selection.errSave': 'Nie udało się zapisać wyboru',
    'selection.errConfirm': 'Nie udało się potwierdzić wyboru',
    'selection.errReopen': 'Nie udało się ponownie otworzyć wyboru zdjęć',
    'selection.saving': 'Zapisywanie…',
    'selection.saveCount': 'Zapisz {count} zdjęć',
    
    // ReadyPanel
    'ready.photosReady': 'zdjęć gotowych',
    'ready.desc': 'Kolejność jest wymieszana i ukryta. Miniaturki pozostają wyłączone do czasu rozpoczęcia każdej rundy rysowania.',
    'ready.hint': 'Oczekiwanie na rozpoczęcie przez administratora…',
    
    // DrawingPanel
    'draw.exitFullscreen': 'Wyjdź z pełnego ekranu',
    'draw.enterFullscreen': 'Pełny ekran',
    'draw.waitingPhoto': 'Oczekiwanie na zdjęcie…',
    'draw.attribution': 'Obraz z',
    'draw.round': 'Runda {current} / {total}',
    'draw.startsIn': 'Runda rozpoczyna się za {count} sekund',
    'draw.remainingAria': 'Pozostały czas rysowania',
    
    // SessionBreakPanels
    'break.takeBreather': 'Złap oddech',
    'break.hiddenDesc': 'Zdjęcie referencyjne jest ukryte do momentu rozpoczęcia następnej rundy.',
    'break.completedRound': 'Ukończono rundę {current} z {total}',
    'break.sessionFinished': 'Sesja zakończona',
    'break.completedRoundsDesc': 'Ukończono {count} rundę. Dobra robota, wszyscy.',
    'break.completedRoundsDescPlural': 'Ukończono {count} rund. Dobra robota, wszyscy.',
    
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
  },
}

export function t(key: string, params?: Record<string, string | number>): string {
  const lang = getBrowserLanguage()
  const template = TRANSLATIONS[lang]?.[key] || TRANSLATIONS.en[key] || key
  if (!params) return template
  return Object.entries(params).reduce(
    (acc, [k, v]) => acc.replace(new RegExp(`{${k}}`, 'g'), String(v)),
    template
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
      description: 'The admin is checking the saved photos before shuffling and starting the session.',
    },
    READY: {
      title: 'Ready to start',
      description: 'Photos are shuffled and hidden. Waiting for the admin to begin.',
    },
    DRAWING: {
      title: 'Draw time',
      description: 'Focus on the reference photo. The timer is server-controlled.',
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
      title: '방장 대기 중',
      description: '방장이 모임을 준비하고 있습니다. 이 페이지에서 잠시 기다려 주세요.',
    },
    SELECTING: {
      title: '레퍼런스 사진 검토 중',
      description: '방장이 셔플 및 세션 시작 전 저장된 사진들을 확인하고 있습니다.',
    },
    READY: {
      title: '시작 준비 완료',
      description: '사진들이 무작위로 섞이고 숨겨졌습니다. 방장이 세션을 시작할 때까지 기다려 주세요.',
    },
    DRAWING: {
      title: '그리기 시간!',
      description: '레퍼런스 사진에 집중해 보세요. 타이머는 서버에 의해 제어됩니다.',
    },
    BETWEEN_ROUNDS: {
      title: '라운드 휴식 시간',
      description: '다음 포즈가 시작되기 전에 짧은 휴식을 취하세요.',
    },
    FINISHED: {
      title: '세션 종료',
      description: '함께 그려주셔서 감사합니다. 다음 주에 또 만나요!',
    },
  },
  ja: {
    WAITING: {
      title: 'ホストの待機中',
      description: '管理者が準備を進めています。このページのままお待ちください。',
    },
    SELECTING: {
      title: 'リファレンス写真の確認中',
      description: '管理者が写真をシャッフルしてセッションを開始する前に、保存された写真を確認しています。',
    },
    READY: {
      title: '開始準備完了',
      description: '写真はシャッフルされ、非表示になっています。管理者が開始するのをお待ちください。',
    },
    DRAWING: {
      title: '描画時間',
      description: 'リファレンス写真に集中してください。タイマーはサーバーによって制御されています。',
    },
    BETWEEN_ROUNDS: {
      title: 'ラウンド間の休憩',
      description: '次のポーズが始まる前に、少し休憩してください。',
    },
    FINISHED: {
      title: 'セッション終了',
      description: '一緒に描いていただきありがとうございました。また来週お会いしましょう。',
    },
  },
  pl: {
    WAITING: {
      title: 'Oczekiwanie na hosta',
      description: 'Administrator przygotowuje pokój. Pozostań na tej stronie.',
    },
    SELECTING: {
      title: 'Przeglądanie zdjęć referencyjnych',
      description: 'Administrator sprawdza zapisane zdjęcia przed ich wymieszaniem i rozpoczęciem sesji.',
    },
    READY: {
      title: 'Gotowy do rozpoczęcia',
      description: 'Zdjęcia są wymieszane i ukryte. Oczekiwanie na rozpoczęcie przez administratora.',
    },
    DRAWING: {
      title: 'Czas na rysowanie',
      description: 'Skoncentruj się na zdjęciu referencyjnym. Licznik czasu jest kontrolowany przez serwer.',
    },
    BETWEEN_ROUNDS: {
      title: 'Przerwa między rundami',
      description: 'Zrób sobie krótką przerwę przed następną pozą.',
    },
    FINISHED: {
      title: 'Sesja zakończona',
      description: 'Dziękujemy za wspólne rysowanie. Do zobaczenia w przyszłym tygodniu.',
    },
  },
}

export function getPhaseMessage(phase: LobbyPhase): PhaseMessage {
  const lang = getBrowserLanguage()
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
      title: '제스처 우선',
      text: '디테일을 무시하고 몸의 흐름과 주요 동세를 먼저 그리세요.',
    },
    {
      title: '눈을 찡그려 단순화하기',
      text: '눈을 약간 감고 보면 자잘한 디테일이 배제되어 명암의 큰 덩어리를 더 쉽게 볼 수 있습니다.',
    },
    {
      title: '팔 전체로 그리기',
      text: '손목만 까딱거리며 짧게 끊어 그리지 말고, 팔꿈치와 어깨를 사용하여 길고 시원한 선을 그리세요.',
    },
    {
      title: '포즈 과장하기',
      text: '짧은 크로키에서는 실제 사진보다 포즈를 더 역동적이고 감정 표현이 풍부하게 그리는 편이 좋습니다.',
    },
    {
      title: '네거티브 스페이스(여백) 활용',
      text: '팔다리 사이에 형성되는 빈 공간의 모양을 관찰하면 비율과 각도를 더 정확하게 잴 수 있습니다.',
    },
    {
      title: '디테일보다 비율 우선',
      text: '몸통과 척추 구조가 확실히 잡힐 때까지는 손, 발, 얼굴 같은 디테일은 비워두세요.',
    },
    {
      title: '유연한 태도 유지',
      text: '이것은 크로키입니다! 완벽하고 다듬어진 그림보다는 속도감과 느낌을 포착하는 것이 훨씬 중요합니다.',
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
}

export function getLocalizedTips(): Tip[] {
  const lang = getBrowserLanguage()
  return LOCALIZED_TIPS[lang] || LOCALIZED_TIPS.en
}
