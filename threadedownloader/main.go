package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

/* ================= CONFIG ================= */

const (
	MAX_THREADS = 16
	USER_AGENT  = "ThreadedDownloader/1.1 (Resume+Tor)"
	TOR_PROXY   = "socks5://127.0.0.1:9050"
	RETRIES     = 3
)

/* ================= COLORS ================= */

const (
	RESET  = "\033[0m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	RED    = "\033[31m"
	BLUE   = "\033[34m"
)

/* ================= UTILS ================= */

func ask(q string) string {
	fmt.Print(q)
	r := bufio.NewReader(os.Stdin)
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}

/* ================= CLIENT / PROXY ================= */

func buildHTTPClient(downloadURL string) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	}

	if strings.Contains(downloadURL, ".onion") {
		fmt.Println(BLUE + "[i] .onion algılandı → TOR aktif" + RESET)
		proxyURL, _ := url.Parse(TOR_PROXY)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

/* ================= PROGRESS ================= */

func printProgress(done, total int64) {
	if total <= 0 {
		return
	}
	percent := float64(done) / float64(total) * 100
	bar := int(percent / 2)
	fmt.Printf("\r[%s%s] %.1f%%",
		strings.Repeat("=", bar),
		strings.Repeat(" ", 50-bar),
		percent,
	)
}

/* ================= RESUME HELPERS ================= */

func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

/* ================= DOWNLOAD PART (RESUME) ================= */

func downloadPart(
	client *http.Client,
	url string,
	start, end int64,
	partPath string,
	progress *int64,
	total int64,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	alreadyDownloaded := fileSize(partPath)
	rangeStart := start + alreadyDownloaded
	if rangeStart >= end {
		*progress += (end - start)
		return
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, end))

	var resp *http.Response
	var err error

	for i := 0; i < RETRIES; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		fmt.Println(RED + "\nParça indirilemedi:" + err.Error() + RESET)
		return
	}
	defer resp.Body.Close()

	f, _ := os.OpenFile(partPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()

	buf := make([]byte, 64*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			f.Write(buf[:n])
			*progress += int64(n)
			printProgress(*progress, total)
		}
		if err == io.EOF {
			break
		}
	}
}

/* ================= MAIN DOWNLOAD ================= */

func download(url string, threads int) {
	client := buildHTTPClient(url)

	req, _ := http.NewRequest("HEAD", url, nil)
	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(RED + "HEAD isteği başarısız" + RESET)
		return
	}
	resp.Body.Close()

	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if size <= 0 {
		fmt.Println(RED + "Dosya boyutu alınamadı" + RESET)
		return
	}

	filename := path.Base(strings.Split(url, "?")[0])
	os.Mkdir("indirilenler", 0755)

	fmt.Println(GREEN+"İndiriliyor:", filename, RESET)

	chunk := size / int64(threads)
	var wg sync.WaitGroup
	var progress int64 = 0

	// Resume için mevcut part dosyalarını hesaba kat
	for i := 0; i < threads; i++ {
		part := fmt.Sprintf("indirilenler/%s.part%d", filename, i)
		progress += fileSize(part)
	}

	for i := 0; i < threads; i++ {
		start := int64(i) * chunk
		end := start + chunk - 1
		if i == threads-1 {
			end = size - 1
		}

		part := fmt.Sprintf("indirilenler/%s.part%d", filename, i)

		wg.Add(1)
		go downloadPart(client, url, start, end, part, &progress, size, &wg)
	}

	wg.Wait()
	fmt.Println()

	// Merge
	finalPath := "indirilenler/" + filename
	out, _ := os.Create(finalPath)

	for i := 0; i < threads; i++ {
		part := fmt.Sprintf("indirilenler/%s.part%d", filename, i)
		in, _ := os.Open(part)
		io.Copy(out, in)
		in.Close()
		os.Remove(part)
	}
	out.Close()

	fmt.Println(GREEN + "✔ İndirme tamamlandı → " + finalPath + RESET)
}

/* ================= MAIN ================= */

func main() {
	fmt.Println(GREEN + "=== Threaded Downloader ===" + RESET)

	url := ask("URL: ")
	if url == "" {
		return
	}

	tStr := ask("Kaç thread? (1-16, öneri:4): ")
	t, _ := strconv.Atoi(tStr)
	if t < 1 {
		t = 4
	}
	if t > MAX_THREADS {
		t = MAX_THREADS
	}

	download(url, t)
}
