package schema

type Signature struct { //batch'in bütünlüğünü taşıyan imza bilgisini içerir
	Algorithm string `json:"algorithm"` //kullanılan imza algoritması bilgisini taşır
	Value     string `json:"value"`     //imza değerinin kendisi
}

type Batch struct {
	BatchID      string    `json:"batch_id"`      //batch'in benzersiz kimliği
	HostID       string    `json:"host_id"`       //batch'ın geldiği sunucunun id'si
	AgentVersion string    `json:"agent_version"` //batch'ı gönderen agent'ın versiyon bilgisi
	CreatedAt    string    `json:"created_at"`    //batch'ın oluşturulduğu an ,RFC3339 formatında string
	Events       []Event   `json:"events"`        //batch içindeki olayların listesi
	Signature    Signature `json:"signature"`     //batch'in bütünlüğünü sağalayan imza bilgisi
}
