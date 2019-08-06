Almost YAML
===========

AYAML, bir YAML düğümünü tek seviyede işleyerek anahtar değer çiftlerinden
oluşan bir eşlem üretmeye yarayan basit ve sınırlı bir YAML biçimidir.

Örnek
-----

```yaml
from /etc/templates/default: debian/buster virtual/lxc class/server   # extension

version:                     0.1.0                                    # string
name:                        ${codename}-${class}_${virtual}          # variable subst
box:                         ${name}_${version}                       # variable subst
github:                      https://github.com/${organization:-omu}  # variable subst
build_id:                    $(date +'%y%m%d%H%M%S')                  # command subst
verbose:                     true                                     # boolean
memory:                      1024                                     # integer
locales:                     [en_US.UTF-8, tr_TR.UTF-8]               # array of strings
preseed:                     ./http/preseed.cfg                       # path subst
pre:                         |                                        # multiline string
                             apt-get update
                             apt-get -y install curl
```

Tanımlar
--------

- Basit anahtar: Olağan bir tanımlayıcı olan, yani `/^[a-zA-Z_][a-zA-Z0-9_]*$/`
  düzenli ifade desenini sağlayan anahtar.

- Skalar değer: Aşağıdaki veri tiplerinden biri.

  + Dizgi
  + Tam sayı
  + Gerçel sayı
  + Mantıksal değer
  + Tarih

- Basit değer: Skalar veya her ögesi skalar olan bir dizi.

- Basit eşlem: Basit anahtar ve basit değer çifti.

- Çıktı: Bir AYAML verisinin ayrıştırılması ve değerlerin işlenmesi sonucunda
  elde edilen ilişkili dizi.

- Çalışma kipi: AYAML ayrıştırıcısının davranış kipleri

  + Gevşek (`lax`) kip: Hatalar göz ardı edilir ve mümkünse çalışmaya devam
    edilir.

  + Sıkı (`strict`) kip: Hatalar göz ardı edilmez; raporlanır ve çalışma hatayla
    sonuçlanır.

  + Güvenli (`secure`) kip: Komut ikameleri göz ardı edilir.

- Literal dizgi: Tek tırnakla ayrılmış veya birden fazla satırdan oluşan dizgi.

- Tırnaksız dizgi: Tırnak kullanılmadan girilen dizgi.

Genel kurallar
--------------

- Her AYAML verisi aynı zamanda geçerli bir YAML verisidir.

- AYAML verisi bir dizi eşlemden oluşur.

- Basit eşlemler değer açıldıktan sonra çıktıda doğrudan yer alır.

- Basit olmayan eşlemler varsa kayıtlı bir eklentiyle tüketilir, yoksa seçilen
  çalışma kipine göre göz ardı edilir veya hata üretir.

İkameler
--------

Literal olmayan dizgi değerlerinde aşağıdaki ikameler gerçekleşir.

### Değişken ikameleri

- `$tanımlayıcı` veya `${tanımlayıcı}` formundaki alt dizgiler ayrıştırma
  tamamlandıktan sonra çıktıda veya proses ortamında bulunan `tanımlayıcı`
  değerle ikame edilir.

- `${tanımlayıcı:-öntanımlı değer}` formunda ikame yapılırken `tanımlayıcı` boş
  dizgiyse `öntanımlı değer` kullanılır.

### Komut ikameleri

- `$(kabuk komutları)` formundaki alt dizgiler ayrıştırma tamamlandıktan sonra
  Bash kabuğuyla çalıştırılarak çıktısı ikame edilir.

- Komut ikamesi güvenli kipte yapılmaz.

### Dosya yolu ikameleri

- Tırnaksız dizgilerde `./path` formundaki sözcükler AYAML verisinin bulunduğu
  dosyaya veya bulunulan dizine göreceli olarak mutlak dosya yoluna çevrilir.

Eklentiler
----------

Eklentiler isteğe bağlıdır.  Bir AYAML gerçeklemesi bir veya birden fazla
eklenti sunabileceği gibi hiç bir eklentiye sahip olmayabilir.

Basit olmayan eşlemlerin yorumlanma semantiği, genel eğilim ayrıştırma hatası
üretmek olmakla birlikte, nihai olarak eklentilerle belirlenir.  Hiç bir eklenti
tanımlı değilse basit olmayan eşlemler `strict` kipte hata üretir, `lax` kipte
eğer değerlendirmeye devam etmek mümkünse göz ardı edilir.

Eklentilere basit olmayan eşlem bilgileriyle birlikte (anahtar ve değer), çıktı
ve bağlam bilgisi iletilir.  Eklenti bu bilgilerden hareketle istediğini
yapmakta özgürdür.  Çıktıda bir değişiklik yapabilir veya yapmayabilir.
