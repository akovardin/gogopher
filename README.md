# gogopher

Gopher-клиент для игровой консоли Playdate, написанный на Go с использованием [pdgo](https://github.com/playdate-go/pdgo).

Позволяет просматривать Gopherspace — меню, файлы и поисковые запросы — прямо на Playdate.

## Сборка

```sh
./build.sh
```

Соберёт `.pdx` для симулятора и устройства.

## Использование

Загрузите `Gopher_sim.pdx` в симулятор Playdate или перенесите `Gopher.pdx` на устройство.

## Зависимости

- Go 1.25
- pdgo v0.8.0
- pdgoc (тулчейн playdate-go)


## Ссылки

bombadillo 192.168.31.92:7070

- хост для тестов 192.168.31.92:7070
- так себе мануал по протоколу https://medium.com/@zoningxtr/ultimate-guide-to-gopher-protocol-from-basics-to-real-exploits-ed2fb788d8e0
- видос "Implementing a Network Protocol in Go" https://www.youtube.com/watch?v=pUaFW98V1Sc
- 2022 - Rocking the Web Bloat: Modern Gopher, Gemini and the Small Internet https://www.youtube.com/watch?v=I2Q35uFCq8Q
- Gemini - The small Internet https://www.ilyameerovich.com/gemini-the-small-internet/
- Introduction to Gemini and the Small Internet https://samsai.eu/post/introduction-to-gemini/
- On Gemini and Gopher, the Internet Protocols https://michaelnordmeyer.com/on-gemini-and-gopher-the-internet-protocols
- https://geminiprotocol.net/
- https://git.sr.ht/~yotam/go-gopher
-  Список килентов: https://en.wikipedia.org/wiki/Gopher_(protocol)#:~:text=Client%20software-,Gopher%20clients,designed%20to%20access%20gopher%20resources.&text=Supports%20page%20cache%2C%20TFTP%20and%20has%20GopherG6%20extension.&text=Eva%20(as%20in%20extra%20vehicular,protocol%20browser%20in%20GTK%204.&text=Supports%20text%20reflow%2C%20bookmarks%2C%20history%2C%20etc.
- browser https://bombadillo.colorfield.space/
- ios browser https://apps.apple.com/at/app/gopher-client/id1235310088
- примеры сайтов https://evertpot.com/
- A Gopher Client in Rust  https://dev.to/krowemoh/notes-on-gopher-266e
- https://ru.wikipedia.org/wiki/Gopher
- https://en.wikipedia.org/wiki/Gemini_(protocol) - очень крутой протокол, который работает с маркдауном
- Почему появился и как устроен протокол Gemini https://habr.com/ru/companies/1cloud/articles/511484/
- Протокол Gemini — минималистичный подход к веб-контенту https://ctf.msk.ru/p/protocol-gemini/
- Хочется странного — шифрование и протокол Gemini https://habr.com/ru/companies/diy_fest/articles/782844/

- Тут делать шрифты https://pdfontconv.frozenfractal.com/
- https://www.fontspace.com/category/pixel,bitmap,cyrillic
- https://www.jetbrains.com/lp/mono/
- https://www.nerdfonts.com/font-downloads
- https://pdfontconv.frozenfractal.com/
- https://idleberg.github.io/playdate-arcade-fonts/

- Как сделать заставку https://www.youtube.com/watch?v=aR-eWv3V-Wo&list=PLOwxD0-Wm6RxpebFlh_-SgcTkBwLTUEN-&index=4

- Gopher браузер https://bombadillo.colorfield.space/