package checks

type Level string

const (
	LevelInfo    Level = "INFO"    //bilgilendirme amaçlı, alarm sayılmıyor
	LevelWarning Level = "WARNING" //kritik eşiğe yaklaşıldığını gösterir
	LevelError   Level = "ERROR"   //kritik eşiğin aşıldığını , gerçekbir alarm durumunu gösterir
)

type CheckResult struct { //bir check'in tek bir ölçümden ürettiği ham sonuç
	CheckType   string  //bu sonucu hangi check'in ürettiğini gösterir
	Level       Level   //sonucun önem derecesini gösterir
	Message     string  //sonucun açıklaması
	MetricName  string  //sonucun ilişkili olduğu ölçümün adı
	MetricValue float64 //ölçülen sayısal değer
	MetricUnit  string  //ölçümün birimi
}

type Check interface { //check interface'i , check tiplerinin implement etmesi gereken metotları tanımlar
	Name() string              //check tipinin adı
	Run() (CheckResult, error) //check tipinin tek bir ölçümden ürettiği ham sonucu döndürür, hata oluşursa error döner
}
