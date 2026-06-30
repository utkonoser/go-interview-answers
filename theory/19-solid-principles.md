# SOLID принципы

## 141. Что такое SOLID принципы и зачем они нужны?

SOLID — пять принципов объектно-ориентированного дизайна для создания поддерживаемого, гибкого и масштабируемого кода. Аббревиатура: S (Single Responsibility), O (Open/Closed), L (Liskov Substitution), I (Interface Segregation), D (Dependency Inversion). Цели: уменьшение связанности (coupling), увеличение связности (cohesion), упрощение тестирования, облегчение изменений и расширения. Применимы не только к ООП — в Go реализуются через интерфейсы, композицию и структуры.

## 142. Single Responsibility Principle (SRP) — что это и как реализуется в Go?

Принцип единственной ответственности: каждый модуль/тип должен иметь только одну причину для изменения, одну зону ответственности. В Go: структура или пакет решает одну задачу. Пример нарушения: структура `User` с методами для валидации, сохранения в БД, отправки email. Правильно: `User` — только данные, `UserValidator` — валидация, `UserRepository` — БД, `EmailService` — отправка. Преимущества: проще тестировать, легче изменять, лучше переиспользовать.

```go
// Плохо: User делает всё
type User struct {
    Name  string
    Email string
}

func (u *User) Validate() error { /* ... */ }
func (u *User) Save() error { /* ... */ }
func (u *User) SendEmail() error { /* ... */ }

// Хорошо: разделение ответственности
type User struct {
    Name  string
    Email string
}

type UserValidator struct{}
func (v *UserValidator) Validate(u *User) error { /* ... */ }

type UserRepository struct{}
func (r *UserRepository) Save(u *User) error { /* ... */ }

type EmailService struct{}
func (s *EmailService) Send(to, subject, body string) error { /* ... */ }
```

## 143. Open/Closed Principle (OCP) — что это и как реализуется в Go?

Принцип открытости/закрытости: код должен быть открыт для расширения, но закрыт для модификации. В Go реализуется через интерфейсы: новый функционал добавляется через новые реализации интерфейса, а не изменением существующего кода. Пример: система обработки платежей — вместо if/switch для каждого типа, интерфейс `PaymentProcessor` с реализациями для разных методов оплаты. Добавление нового метода не требует изменения существующего кода.

```go
// Плохо: нужно модифицировать при добавлении метода
func ProcessPayment(method string, amount float64) error {
    switch method {
    case "card":
        return processCard(amount)
    case "paypal":
        return processPaypal(amount)
    // при добавлении нового метода нужно менять этот код
    }
}

// Хорошо: открыто для расширения через интерфейс
type PaymentProcessor interface {
    Process(amount float64) error
}

type CardProcessor struct{}
func (p *CardProcessor) Process(amount float64) error { /* ... */ }

type PaypalProcessor struct{}
func (p *PaypalProcessor) Process(amount float64) error { /* ... */ }

type CryptoProcessor struct{} // новый тип без изменения существующего кода
func (p *CryptoProcessor) Process(amount float64) error { /* ... */ }

func ProcessPayment(processor PaymentProcessor, amount float64) error {
    return processor.Process(amount)
}
```

## 144. Liskov Substitution Principle (LSP) — что это и как реализуется в Go?

Принцип подстановки Барбары Лисков: объекты должны быть заменяемыми на экземпляры их подтипов без изменения правильности программы. В Go: любая реализация интерфейса должна выполнять контракт интерфейса корректно. Нарушение: если метод паникует, возвращает неожиданные значения или требует специальной обработки для конкретной реализации. Правильно: все реализации ведут себя предсказуемо согласно контракту интерфейса.

```go
// Интерфейс для хранилища
type Storage interface {
    Save(key string, value []byte) error
    Load(key string) ([]byte, error)
}

// Хорошо: обе реализации корректно выполняют контракт
type FileStorage struct{}
func (s *FileStorage) Save(key string, value []byte) error {
    return os.WriteFile(key, value, 0644)
}
func (s *FileStorage) Load(key string) ([]byte, error) {
    return os.ReadFile(key)
}

type MemoryStorage struct {
    data map[string][]byte
}
func (s *MemoryStorage) Save(key string, value []byte) error {
    s.data[key] = value
    return nil
}
func (s *MemoryStorage) Load(key string) ([]byte, error) {
    val, ok := s.data[key]
    if !ok {
        return nil, os.ErrNotExist
    }
    return val, nil
}

// Плохо: ReadOnlyStorage нарушает контракт (Save всегда падает)
type ReadOnlyStorage struct{}
func (s *ReadOnlyStorage) Save(key string, value []byte) error {
    return errors.New("read-only storage") // нарушение LSP
}
func (s *ReadOnlyStorage) Load(key string) ([]byte, error) {
    return os.ReadFile(key)
}

// Правильно: отдельный интерфейс для read-only
type ReadOnlyStorageInterface interface {
    Load(key string) ([]byte, error)
}
```

## 145. Interface Segregation Principle (ISP) — что это и как реализуется в Go?

Принцип разделения интерфейсов: клиенты не должны зависеть от методов, которые они не используют. В Go: создавать маленькие, специфичные интерфейсы вместо больших универсальных. Go идиома: "accept interfaces, return structs" и "чем меньше интерфейс, тем он полезнее". Пример: вместо одного `Database` интерфейса с 20 методами — несколько маленьких: `Reader`, `Writer`, `Transactor`. Стандартная библиотека Go следует ISP: `io.Reader`, `io.Writer`, `io.Closer`.

```go
// Плохо: большой интерфейс, клиенты зависят от ненужных методов
type Database interface {
    Connect() error
    Disconnect() error
    Query(sql string) ([]Row, error)
    Execute(sql string) error
    BeginTransaction() error
    Commit() error
    Rollback() error
    Backup() error
    Restore() error
}

// Функция нуждается только в чтении, но зависит от всего интерфейса
func GetUsers(db Database) ([]User, error) {
    rows, err := db.Query("SELECT * FROM users")
    // ...
}

// Хорошо: разделенные интерфейсы
type Querier interface {
    Query(sql string) ([]Row, error)
}

type Executor interface {
    Execute(sql string) error
}

type Transactor interface {
    BeginTransaction() error
    Commit() error
    Rollback() error
}

// Функция зависит только от того, что использует
func GetUsers(q Querier) ([]User, error) {
    rows, err := q.Query("SELECT * FROM users")
    // ...
}

// Композиция интерфейсов для комплексных операций
type Database interface {
    Querier
    Executor
    Transactor
}

// Примеры из стандартной библиотеки Go (ISP в действии)
// io.Reader, io.Writer, io.Closer - маленькие, специфичные интерфейсы
func ProcessData(r io.Reader) error {
    // нужен только Reader, не зависит от Writer или Closer
}
```

## 146. Dependency Inversion Principle (DIP) — что это и как реализуется в Go?

Принцип инверсии зависимостей: модули высокого уровня не должны зависеть от модулей низкого уровня, оба должны зависеть от абстракций. В Go: зависимость от интерфейсов, а не конкретных реализаций. Dependency Injection через конструкторы или параметры функций. Пример: сервис зависит от интерфейса `Repository`, а не конкретной `PostgresRepository`. Это позволяет легко менять реализацию и мокать при тестировании.

```go
// Плохо: зависимость от конкретной реализации
type UserService struct {
    db *PostgresDatabase // зависит от конкретной реализации
}

func NewUserService() *UserService {
    return &UserService{
        db: &PostgresDatabase{}, // жесткая связь
    }
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.db.FindUserByID(id)
}

// Хорошо: зависимость от абстракции (интерфейса)
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}

type UserService struct {
    repo UserRepository // зависит от интерфейса
}

// Dependency Injection через конструктор
func NewUserService(repo UserRepository) *UserService {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.FindByID(id)
}

// Конкретные реализации
type PostgresRepository struct {
    db *sql.DB
}

func (r *PostgresRepository) FindByID(id int) (*User, error) {
    // PostgreSQL implementation
}

type MongoRepository struct {
    client *mongo.Client
}

func (r *MongoRepository) FindByID(id int) (*User, error) {
    // MongoDB implementation
}

// Легко менять реализацию
func main() {
    // Production
    postgresRepo := &PostgresRepository{db: db}
    service := NewUserService(postgresRepo)
    
    // Testing
    mockRepo := &MockRepository{}
    testService := NewUserService(mockRepo)
}

// Mock для тестирования
type MockRepository struct {
    users map[int]*User
}

func (r *MockRepository) FindByID(id int) (*User, error) {
    if user, ok := r.users[id]; ok {
        return user, nil
    }
    return nil, errors.New("not found")
}
```

## 147. Как SOLID принципы помогают в тестировании Go кода?

SOLID принципы делают код легко тестируемым: 1) SRP — изолированные компоненты легко тестировать по отдельности; 2) OCP — новые тесты без изменения существующих; 3) LSP — мокать любую реализацию интерфейса; 4) ISP — мокать только нужные методы, минимальные зависимости; 5) DIP — легко подменять зависимости на моки через интерфейсы. Dependency Injection позволяет передавать тестовые реализации. В Go: моки через интерфейсы, table-driven tests, пакеты testify/mock или gomock для генерации моков.

```go
// Благодаря DIP и ISP легко тестировать
type UserService struct {
    repo UserRepository
    emailService EmailSender // маленький интерфейс
}

type EmailSender interface {
    Send(to, subject, body string) error
}

func (s *UserService) RegisterUser(email string) error {
    user := &User{Email: email}
    if err := s.repo.Save(user); err != nil {
        return err
    }
    return s.emailService.Send(email, "Welcome", "Thanks for registering")
}

// Тестирование с моками
func TestRegisterUser(t *testing.T) {
    // Моки реализуют те же интерфейсы
    mockRepo := &MockRepository{
        users: make(map[int]*User),
    }
    mockEmail := &MockEmailSender{
        sentEmails: []string{},
    }
    
    service := NewUserService(mockRepo, mockEmail)
    
    err := service.RegisterUser("test@example.com")
    assert.NoError(t, err)
    assert.Equal(t, 1, len(mockEmail.sentEmails))
}

type MockEmailSender struct {
    sentEmails []string
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.sentEmails = append(m.sentEmails, to)
    return nil
}
```

## 148. Какие компромиссы и проблемы могут возникать при применении SOLID?

Чрезмерное применение SOLID может привести к: 1) Over-engineering — слишком много мелких интерфейсов и типов для простых задач; 2) Снижение производительности — дополнительные уровни абстракции, indirect calls через интерфейсы; 3) Усложнение кода — труднее следить за потоком выполнения; 4) Преждевременная абстракция — создание интерфейсов "на будущее". Баланс: начинать с простого решения, рефакторить при появлении реальной необходимости. "Make it work, make it right, make it fast". В Go: избегать интерфейсов пока не нужны две реализации или моки для тестов.

```go
// Over-engineering: слишком много абстракций для простой задачи
type ConfigReader interface {
    Read() (Config, error)
}

type ConfigValidator interface {
    Validate(Config) error
}

type ConfigLoader interface {
    Load() (Config, error)
}

type ConfigLoaderImpl struct {
    reader    ConfigReader
    validator ConfigValidator
}

// Для простого случая достаточно одной функции
func LoadConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, err
    }
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, err
    }
    return cfg, nil
}

// Правило: "Don't create abstractions, discover them"
// Создавать интерфейсы когда появляется вторая реализация или нужны моки
```


