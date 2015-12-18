package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"image"
	"image/jpeg"
	"image/png"
	"workshop/blog"
	"workshop/storage"

	"github.com/nfnt/resize"
)

// Storage is a structure of blog post
type Storage struct {
	ID       int
	Title    string
	Body     string
	CreateAt string
	Picture  string
}

// BlogStorage will store blog posts
var BlogStorage storage.Storage

// render is a function to render template
func render(w http.ResponseWriter, sTemplateName string, data map[string]interface{}) {

	// all templates are in a view folder.
	// we reference them by their names
	sNewTemplateName := "./view/" + sTemplateName + ".tpl"
	t, err := template.New("layout").ParseFiles("./view/layout.tpl", sNewTemplateName)

	if sTemplateName == "edit" || sTemplateName == "post" {
		t, err = t.ParseFiles("./view/form.tpl")
	}

	if err != nil {
		fmt.Println(w, "Error\n", err.Error())
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Println(w, "Error\n", err.Error())
	}

}

// isError is a function to handle errors
func isError(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return true
	}
	return false
}

// handlePost is a method which handle requests
// associated with blog entry creation
func handlerPost(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		data := make(map[string]interface{})
		data["PageTitle"] = "Create new post"
		data["action"] = "create"
		render(w, "post", data)

	case "POST":

		switch r.FormValue("action") {
		case "create":
			createAt := time.Now().Format(time.RFC822) // time format into string
			id := BlogStorage.GetLength() + 1
			title := r.FormValue("title")
			body := r.FormValue("body")

			entry := blog.Entry{ID: id, Title: title, Body: body, CreateAt: createAt}
			entryID := BlogStorage.Add(entry) // we created entry, now we can save it

			r.ParseMultipartForm(4096 * 10) // how many bytes of image will be stored in memory
			f, fh, err := r.FormFile("picture")
			if err != nil {
				log.Println(BlogStorage.GetAll())
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			ext := ""
			switch fh.Header.Get("Content-Type") {
			case "image/png":
				ext = "png"
			case "image/jpeg":
				ext = "jpeg"
			}

			// here we are creating random file name
			randBytes := make([]byte, 16)
			rand.Read(randBytes)
			filename := filepath.Join(hex.EncodeToString(randBytes)) + "." + ext
			fdesc, err := os.Create("./img/" + filename)
			if isError(err, w) {
				return
			}

			// after successfull creation we copy our image to this file
			io.Copy(fdesc, f)
			fdesc.Close()
			f.Close()

			// after successfull save of image to file,
			// we update our created entry
			entry.Picture = filename
			BlogStorage.Update(entryID, entry)

			sImageType := ""
			switch fh.Header.Get("Content-Type") {
			case "image/png":
				sImageType = "png"
			case "image/jpeg":
				sImageType = "jpeg"
			}
			
			// here we are making three files
			// of different sizes from original image
			for i := 2; i <= 3; i++ {
				// here we are using goroutines to do resize operation concurrently
				go func(divider int, sOriginalFile string, sType string, w http.ResponseWriter) {

					file, err := os.Open("./img/" + sOriginalFile)
					if isError(err, w) {
						return
					}

					// decode jpeg into image.Image
					var (
						img  image.Image
						conf image.Config
					)
					switch sType {
					case "jpeg":
						// when we decoded config, our bytes shifted
						conf, err = jpeg.DecodeConfig(file)
						if isError(err, w) {
							return
						}
						// so we need to return them back
						file.Seek(0, 0)
						img, err = jpeg.Decode(file)
						if isError(err, w) {
							return
						}
					case "png":
						img, err = png.Decode(file)
						if isError(err, w) {
							return
						}
						file.Seek(0, 0)
						conf, err = png.DecodeConfig(file)
						if isError(err, w) {
							return
						}

					}
					file.Close()

					// resize to width / divider using Lanczos resampling
					// and preserve aspect ratio
					m := resize.Resize(uint(conf.Width/divider), 0, img, resize.Lanczos3)

					out, err := os.Create("./img/" + strconv.Itoa(divider) + "_" + sOriginalFile)
					if isError(err, w) {
						return
					}
					defer out.Close()

					// write new image to file
					switch sType {
					case "jpeg":
						jpeg.Encode(out, m, nil)
					case "png":
						png.Encode(out, m)
					}
				}(i, filename, sImageType, w)
			}

		case "edit":
			var (
				err error
				id  int
			)

			if id, err = strconv.Atoi(r.FormValue("id")); err != nil {
				log.Println(err)
				http.Error(w, "ID not found", 404)
				return
			}

			entry, err := BlogStorage.GetByID(id)
			if err != nil {
				http.Error(w, "ID not found", 404)
				return
			}

			entry.Body = r.FormValue("body")
			entry.Title = r.FormValue("title")

			BlogStorage.Update(entry.ID, entry)

		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

}

// handlerEdit handle blog entry edit request
func handlerEdit(w http.ResponseWriter, r *http.Request) {

	validPath := regexp.MustCompile("^/(edit)/([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if len(m) < 2 {
		http.NotFound(w, r)
		return
	}

	id := 0
	var err error
	if id, err = strconv.Atoi(m[2]); err != nil {
		http.NotFound(w, r)
		return
	}

	entry, err := BlogStorage.GetByID(id)
	if isError(err, w) {
		http.NotFound(w, r)
		return

	}

	data := make(map[string]interface{})
	data["PageTitle"] = "Edit"
	data["title"] = entry.Title
	data["body"] = entry.Body
	data["id"] = id
	data["action"] = "edit"
	render(w, "post", data)
}

func main() {

	log.Println("We ready to start")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		data["PageTitle"] = "Workshop blog"
		data["storage"] = BlogStorage.GetAll()
		render(w, "main", data)
	})

	http.HandleFunc("/edit/", handlerEdit)
	http.HandleFunc("/post", handlerPost)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))

	log.Fatal(http.ListenAndServe(":9000", nil))
}
