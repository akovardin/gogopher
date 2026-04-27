# gogopher

Gopher-клиент для игровой консоли Playdate, написанный на Go с использованием [pdgo](https://github.com/playdate-go/pdgo).

Позволяет просматривать Gopherspace — меню и файлы прямо на Playdate.

## Сборка и использование

Для сборки версии под симулятор используйте `task build` и после этого можно перетягивать файл `Gopher_sim.pdx` в симулятор

Чтобы запустить на устройстве - выполните команду `task deploy`


## Зависимости

- Go 1.25
- pdgo v0.8.0
- pdgoc (тулчейн playdate-go)
- [taskfile.dev](https://taskfile.dev/)
