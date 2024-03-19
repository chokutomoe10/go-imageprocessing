package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"gocv.io/x/gocv"
)

func main() {
	http.HandleFunc("/convert", Convert)
	http.HandleFunc("/resize", Resize)
	http.HandleFunc("/compress", Compress)

	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func Convert(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["Images"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileByte, errByte := ioutil.ReadAll(file)
		if errByte != nil {
			http.Error(w, errByte.Error(), http.StatusInternalServerError)
			return
		}

		fileType := http.DetectContentType(fileByte)

		if fileType == "image/png" {
			tempFilePng, err := ioutil.TempFile("images/converted/PNG", "image-*.png")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer tempFilePng.Close()

			tempFilePng.Write(fileByte)

			img := gocv.IMRead(tempFilePng.Name(), gocv.IMReadColor)
			if img.Empty() {
				http.Error(w, "Couldn't read the image", http.StatusInternalServerError)
				return
			}
			defer img.Close()

			convertedImage, errConvert := gocv.IMEncode(".jpg", img)
			if errConvert != nil {
				http.Error(w, errConvert.Error(), http.StatusInternalServerError)
				return
			}

			tempFileJpeg, errJpeg := ioutil.TempFile("images/converted/JPEG", "image-*.jpeg")
			if errJpeg != nil {
				http.Error(w, errJpeg.Error(), http.StatusInternalServerError)
				return
			}
			defer tempFileJpeg.Close()

			tempFileJpeg.Write(convertedImage.GetBytes())
		} else {
			http.Error(w, "Image file must be PNG", http.StatusBadRequest)
			return
		}
	}

	fmt.Fprint(w, "Successfully Converted Files")
}

func Resize(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errForm := r.ParseForm(); errForm != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["Images"]

	sizeX, _ := strconv.Atoi(r.PostForm.Get("SizeX"))
	sizeY, _ := strconv.Atoi(r.PostForm.Get("SizeY"))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileByte, errByte := ioutil.ReadAll(file)
		if errByte != nil {
			http.Error(w, errByte.Error(), http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(fileHeader.Filename)

		tempFile, errFile := ioutil.TempFile("images/uploads", "image-*"+ext)
		if errFile != nil {
			http.Error(w, errFile.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		tempFile.Write(fileByte)

		img := gocv.IMRead(tempFile.Name(), gocv.IMReadColor)
		if img.Empty() {
			http.Error(w, "Error: Couldn't read the image", http.StatusInternalServerError)
			return
		}
		defer img.Close()

		resized := gocv.NewMat()
		gocv.Resize(img, &resized, image.Pt(sizeX, sizeY), 0, 0, gocv.InterpolationDefault)

		tempResizedFile, errResize := ioutil.TempFile("images/resized", "image-*"+ext)
		if errResize != nil {
			http.Error(w, errResize.Error(), http.StatusInternalServerError)
			return
		}
		defer tempResizedFile.Close()

		switch ext {
		case ".jpg", ".jpeg":
			convertedImage, err := gocv.IMEncode(".jpg", resized)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tempResizedFile.Write(convertedImage.GetBytes())
			continue
		case ".png":
			convertedImage, err := gocv.IMEncode(".png", resized)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tempResizedFile.Write(convertedImage.GetBytes())
			continue
		default:
			fmt.Println("Unsupported file type:", ext)
			continue
		}
	}

	fmt.Fprint(w, "Successfully Resized Files")
}

func Compress(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errForm := r.ParseForm(); errForm != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["Images"]

	sizeX, _ := strconv.Atoi(r.PostForm.Get("SizeX"))
	sizeY, _ := strconv.Atoi(r.PostForm.Get("SizeY"))
	sizeJPEG, _ := strconv.Atoi(r.PostForm.Get("SizeJPEG"))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileByte, errByte := ioutil.ReadAll(file)
		if errByte != nil {
			http.Error(w, errByte.Error(), http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(fileHeader.Filename)

		tempFile, errFile := ioutil.TempFile("images/uploads", "image-*"+ext)
		if errFile != nil {
			http.Error(w, errFile.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		tempFile.Write(fileByte)

		img := gocv.IMRead(tempFile.Name(), gocv.IMReadColor)
		if img.Empty() {
			http.Error(w, "Error: Couldn't read the image", http.StatusInternalServerError)
			return
		}
		defer img.Close()

		tempCompressedFile, errCompress := ioutil.TempFile("images/compressed", "image-*"+ext)
		if errCompress != nil {
			http.Error(w, errCompress.Error(), http.StatusInternalServerError)
			return
		}
		defer tempCompressedFile.Close()

		switch ext {
		case ".jpg", ".jpeg":
			compressedImage, err := gocv.IMEncodeWithParams(".jpg", img, []int{gocv.IMWriteJpegQuality, sizeJPEG})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tempCompressedFile.Write(compressedImage.GetBytes())
			continue
		case ".png":
			resized := gocv.NewMat()
			gocv.Resize(img, &resized, image.Pt(sizeX, sizeY), 0, 0, gocv.InterpolationDefault)

			convertedImage, err := gocv.IMEncode(".png", resized)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tempCompressedFile.Write(convertedImage.GetBytes())
			continue
		default:
			fmt.Println("Unsupported file type:", ext)
			continue
		}
	}

	fmt.Fprint(w, "Successfully Compressed Files")
}
