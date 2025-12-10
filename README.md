# Threaded Downloader

Threaded Downloader, Go dili ile yazılmış; çok parçalı indirme yapabilen, otomatik kaldığı yerden devam etme özelliğine sahip bir komut satırı  aracıdır..

Bu proje; büyük dosyaların indirilmesi, sık kopan internet bağlantıları, Tor ağı üzerinden anonim indirme gibi özellikler göz önünde bulundurularak tasarlanmıştır.

# Özellikler

• Çok parçalı  indirme desteği (1–16 thread)

• Otomatik resume (indirmenin yarım kalması durumunda kaldığı yerden devam eder)

• Dahili progress bar 

• Özel User-Agent kullanımı

• Proxy desteği (HTTP / HTTPS / SOCKS5)

• .onion uzantılı URL’ler için otomatik Tor kullanımı

• Hata durumlarında yeniden deneme  mekanizması

• Timeout yönetimi

• Renkli terminal çıktıları

# Proxy ve Tor Desteği

Program, girilen URL’ye göre bağlantı yöntemini otomatik belirler:

• Normal web URL’leri → Doğrudan internet bağlantısı

• .onion uzantılı URL’ler → Otomatik olarak Tor ağı üzerinden indirme

Tor için kullanılan proxy adresi:

socks5://127.0.0.1:9050

.onion adreslerinden indirme yapılabilmesi için Tor Browser veya Tor servisinin çalışır durumda olması gerekir.


# Otomatik Resume

İndirme işlemi şu durumlarda yarıda kesilse bile veri kaybı yaşanmaz:

• İnternet bağlantısı koparsa

• Program kapanırsa

• Kullanıcı Ctrl + C ile işlemi durdurursa

Program tekrar çalıştırıldığında:

Aynı URL girilir

Mevcut .part dosyaları kontrol edilir

İndirme, kaldığı byte’tan otomatik olarak devam eder

Resume işlemi için ek .meta veya yapılandırma dosyası kullanılmaz.


# Gereksinimler

• Go 1.20 veya üzeri

• (Opsiyonel) Tor Browser veya Tor servisi


# Kurulum ve Çalıştırma

Proje dizinine girip aşağıdaki komutu çalıştırmak yeterlidir:

go run main.go

Program çalıştığında kullanıcıdan sırasıyla şu bilgiler istenir:

URL:
Kaç thread? (1–16, öneri: 4)


# İndirilen Dosyalar

İndirilen tüm dosyalar otomatik olarak aşağıdaki klasöre kaydedilir:

indirilenler/

İndirme tamamlandıktan sonra geçici .part dosyaları otomatik olarak silinir.

# Örnek Çıktı

=== Threaded Downloader ===

URL: https://example.com/largefile.iso

Kaç thread? (1–16, öneri:4): 8

[========================= ] 54.4%
✔ İndirme tamamlandı → indirilenler/largefile.iso



Projeyi faydalı bulduysanız GitHub üzerinde yıldız vererek destek olabilirsiniz.
