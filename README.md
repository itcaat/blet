# blet

CLI тул для поиска дешевых авиабилетов в консоли.

```
/cmd
  ├── root.go              // точка входа CLI
  ├── cheapest.go          // CLI-обвязка для дешевых билетов
  └── week.go              // CLI-обвязка для билетов на неделю

/internal/usecase
  ├── cheapest.go          // usecase: логика получения дешевых билетов
  └── week.go              // usecase: логика для week matrix

/internal/api
  └── tpclient.go          // взаимодействие с API TravelPayouts

/internal/form
  └── citypairs.go         // интерактивные формы на huh

/internal/models
  └── ticket.go

/internal/cache
  └── cache.go

/internal/cmd
  └── helpers.go

/config
  └── config.go

/main.go

```
