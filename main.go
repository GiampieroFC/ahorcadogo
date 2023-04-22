package main

import (
	// "errors"
	"fmt"
	"github/GiampieroFC/ahorcadoGO/scrap"
	"github/GiampieroFC/ahorcadoGO/widgets"
	"net/url"
	"strings"

	// "strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var re bool = false
var u url.URL
var pal_u string

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("AhorcadoGO")

	channel := make(chan string)

	letrasEscritas := binding.NewString()
	var le string
	pantallaBinding := binding.NewString()
	chances := binding.NewString()
	pantallaBinding.Set("Buscando artículo en Wikipedia...")
	chances.Set("mmm...")

	letrasEscritasLabel := widget.NewLabelWithData(letrasEscritas)
	letrasEscritasLabel.Hide()

	check := widget.NewCheck(
		"Ver letras ya escritas: ",
		func(value bool) {
			fmt.Println("Radio set to", value)
			if value {
				letrasEscritas.Set("(carácteres especiales ya escritos) " + le)
				letrasEscritasLabel.Show()
			} else {
				letrasEscritasLabel.Hide()
			}
		})
	check.MinSize()

	label := widget.NewLabelWithData(pantallaBinding)
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle.Monospace = true

	opor := widget.NewLabelWithData(chances)
	opor.Wrapping = fyne.TextWrapWord
	opor.Alignment = fyne.TextAlignCenter

	link := widget.NewHyperlink(pal_u, &u)
	link.Hide()

	input := widgets.CrearEntradaConEvento()
	input.Cursor()

	enter := func(key *fyne.KeyEvent) {
		fmt.Println("has presionado", key.Name)
		if input.Text == "" {
			return
		}
		if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
			if !strings.Contains(le, strings.ToLower(input.Text)) {
				le += strings.ToLower(input.Text) + " "
				letrasEscritas.Set(le)
			}
			channel <- strings.ToLower(input.Text)
			input.SetText("")
		}
	}
	input.TeclaEvento = enter

	input.Validator = func(s string) error {
		fmt.Printf("len(sa): %v\n", len(s))
		if len(s) > 1 && s != "ñ" {
			input.SetText("")
			input.SetPlaceHolder("ESCRIBE SOLO UNA LETRA")
			return nil
		}
		input.SetPlaceHolder("escribe una letra")
		return nil
	}

	botonE := widget.NewButton("Enter", func() {
		if input.Text == "" {
			return
		}
		if !strings.Contains(le, strings.ToLower(input.Text)) {
			le += strings.ToLower(input.Text) + " "
			letrasEscritas.Set(le)
		}
		channel <- strings.ToLower(input.Text)
		input.SetText("")
		fmt.Println("has escrito", strings.ToLower(input.Text))

	})

	logica(pantallaBinding, chances, link, input, botonE, channel)

	Iniciar := func() {
		if !input.Disabled() && !botonE.Disabled() {
			re = true
			channel <- strings.ToLower(input.Text)
			input.SetText("")
			pantallaBinding.Set("")
		}
		le = ""
		letrasEscritas.Set("")
		logica(pantallaBinding, chances, link, input, botonE, channel)
	}

	botonR := widget.NewButton("Reiniciar", Iniciar)

	header := container.New(layout.NewHBoxLayout(), check, letrasEscritasLabel, layout.NewSpacer(), botonR)

	separador := widget.NewSeparator()

	content := container.New(layout.NewVBoxLayout(), header, separador, opor, label, link, input, botonE)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 300))
	myWindow.ShowAndRun()
}

func logica(pantallaBinding binding.String, chances binding.String, link *widget.Hyperlink, input *widgets.EntradaConEvento, botonE *widget.Button, channel chan string) {
	palabra := scrap.PrintPalabra()
	var letras string = " ()-,–.áàäãâāæéèëêēíìïīóòöōøúùüūßšçπ0123456789€%$&#'\"\\/[]{}!¡¿?+=:_<>·@#|°®"
	if link.Visible() {
		link.Hide()
	}
	if input.Disabled() && botonE.Disabled() {
		input.Enable()
		botonE.Enable()
	}

	re = false
	go loop(pantallaBinding, chances, link, input, botonE, channel, palabra, letras)

}

func loop(pantallaBinding binding.String, chances binding.String, link *widget.Hyperlink, input *widgets.EntradaConEvento, botonE *widget.Button, channel chan string, palabra scrap.Paldef, letras string) {
	var l string
	var Ganaste bool
	lowerPalabra := strings.ToLower(palabra.Palabra)
	oportunidades := len(lowerPalabra) + 2

	for {
		fmt.Printf("re: %v\n", re)
		if re {
			pantallaBinding.Set("Salimos!")
			return
		}
		Ganaste = true
		setP := make(map[string]string)
		for _, v := range lowerPalabra {
			s := string(v)
			setP[s] = string(v)
		}

		setL := make(map[string]string)
		for _, v := range letras {
			s := string(v)
			setL[s] = string(v)
		}
		fmt.Println(setL)

		var pantalla string
		for _, v := range lowerPalabra {
			e, encontrado := setL[string(v)]
			if e == " " {
				pantalla += "  "
			}
			if encontrado {
				pantalla += string(v) + " "
			} else {
				pantalla += "_ "
			}
		}
		pantallaBinding.Set(pantalla)

		cont_chances := fmt.Sprintf("tienes %d oportunidades", oportunidades)

		chances.Set(cont_chances)

		for v := range setP {
			_, encontrado := setL[v]
			if !encontrado {
				Ganaste = false
				break
			}
		}
		if Ganaste {
			ganaste := fmt.Sprintf("GANASTE!\n\n%s\n\n", palabra.Definicion)
			u = *palabra.Link
			link.Text = palabra.Palabra
			link.Show()
			fmt.Println(ganaste)
			pantallaBinding.Set(ganaste)
			input.Disable()
			botonE.Disable()
			return
		} else {
			fmt.Println("\nAÚN NO GANAS")
		}
		if oportunidades == 0 {
			perdiste := fmt.Sprintf("¡PERDISTE!\n\n%s\n\n", palabra.Definicion)
			u = *palabra.Link
			link.Text = palabra.Palabra
			link.Show()
			fmt.Println(perdiste)
			pantallaBinding.Set(perdiste)
			input.Disable()
			botonE.Disable()
			return
		}

		l = <-channel
		_, hay := setP[l]
		if !hay {
			oportunidades--
		}
		letras += l
		fmt.Printf("letras: %v\n", letras)
	}

}
