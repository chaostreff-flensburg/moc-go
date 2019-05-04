# MOC Go Client

[Message Operation Center](https://github.com/chaostreff-flensburg/moc) Go Client

## Usage

```bash
import mocgo "github.com/chaostreff-flensburg/moc-go"
```

Get Messagelist:

```golang
client := mocgo.NewClient("http://localhost:8080")
messages := client.Request()
```

Always give the latest Message by a channel:

```golang
client := mocgo.NewClient("http://localhost:8080")

client.Loop(20 * time.Second)

for message := range client.NewMessages {
	fmt.Println(message.ID)
}
```
