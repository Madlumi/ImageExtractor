package main

import (
   "bufio"
   "fmt"
   "net/http"
   "os"
   "os/exec"
   "strings"
   "io"
   "github.com/PuerkitoBio/goquery"
)

func E(err error) bool { return err != nil }
func isImageURL(url string) bool {
   return strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".gif")
}

func extractImageLinksFromURL(url string) ([]string, error) {
   res, err := http.Get(url)
   if E(err) {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("Failed to fetch the webpage. Status code: %d", res.StatusCode)
   }
   return extractImageLinksFromReader(res.Body)
}

func extractImageLinksFromFile(path string) ([]string, error) {
   file, err := os.Open(path)
   if E(err) {
      return nil, err
   }
   defer file.Close()
   return extractImageLinksFromReader(file)
}

func extractImageLinksFromReader(reader io.Reader) ([]string, error) {
   doc, err := goquery.NewDocumentFromReader(reader)
   if E(err) {
      return nil, err
   }
   var imageLinks []string
   doc.Find("img, a").Each(func(i int, element *goquery.Selection) {
      if element.Is("img") {
         src, exists := element.Attr("src")
         if exists && src != "" {
            imageLinks = append(imageLinks, src)
         }
      } else if element.Is("a") {
         href, exists := element.Attr("href")
         if exists && isImageURL(href) {
            imageLinks = append(imageLinks, href)
         }
      }
   })
   return imageLinks, nil
}


func blacklist(l string) bool {
   banned := []string{
      "thumbnail",
      "svg",
      "-logo",
      "extras/store",
      "/search",
      "trace.moe", 
      "iqdb.org",
   }

   for _, ban := range banned {
      if strings.Contains(l, ban) {
         return true
      }
   }
   return false
}

func openLink(image string) {
   cmd := exec.Command("xdg-open", image)
   err := cmd.Run()
   if err != nil { fmt.Println("An error occurred while opening URLs:", err); }
}
func main() {
   if len(os.Args) > 1 {
      for _, arg := range os.Args[1:] {
         processInput(arg)
      }
   } else {
      scanner := bufio.NewScanner(os.Stdin)
      for scanner.Scan() {
         input := scanner.Text()
         processInput(input)
      }

      if err := scanner.Err(); err != nil {
         fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
      }
   }
}

func processInput(input string) {
   var images []string;
   var err error;
   if (strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")) {
      images, err = extractImageLinksFromURL(input)
   } else {
      images, err = extractImageLinksFromFile(input)
   }
   if E(err) { fmt.Println("Error:", err); return; }
   if len(images) > 0 {
      fmt.Println("Image links found:")
      uniq := make(map[string]bool)
      for _, img := range images {
         if blacklist(img){continue; }
         if(uniq[getname(img)]){continue;}
         uniq[getname(img)]=true;
         fmt.Println(img);
         openLink(img)
      }
   }
}
func getname(s string)string{
    lastIndex := strings.LastIndex(s, "/")
    if lastIndex != -1 && lastIndex < len(s)-1 {
        return s[lastIndex+1:]
    }
    return s
}


































