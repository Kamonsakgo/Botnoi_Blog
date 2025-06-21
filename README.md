# Botnoi Blog API

## คำอธิบายโปรเจ็ค
API สำหรับจัดการบล็อกและกิจกรรมไฮไลท์ (Highlight Events) พัฒนาด้วย Go Fiber พร้อมการเชื่อมต่อฐานข้อมูล MongoDB และระบบจัดเก็บไฟล์ด้วย AWS S3

## เทคโนโลยีที่ใช้
- **Go Fiber** - Web Framework สำหรับ Go
- **MongoDB** - ฐานข้อมูล NoSQL
- **AWS S3** - บริการจัดเก็บไฟล์
- **JWT** - การจัดการ Authentication
- **Swagger** - API Documentation

## โครงสร้างโปรเจ็ค

```
Botnoi_Blog/
├── configuration/        # การตั้งค่า Fiber
├── domain/              # โดเมนของโปรเจ็ค
│   ├── datasources/     # การเชื่อมต่อฐานข้อมูล
│   ├── entities/        # โครงสร้างข้อมูล
│   └── repositories/    # การเข้าถึงข้อมูล
├── src/
│   ├── gateways/       # HTTP Handlers
│   ├── middlewares/    # Middleware functions
│   ├── services/       # Business Logic
│   └── utils/          # ฟังก์ชันช่วยเหลือ
├── tests/              # ไฟล์ทดสอบ
├── docs/               # API Documentation
└── main.go             # ไฟล์หลัก
```

## การติดตั้งและเรียกใช้งาน

### ข้อกำหนดเบื้องต้น
- Go 1.19 หรือสูงกว่า
- MongoDB
- AWS S3 Account

### การติดตั้ง

1. **Clone โปรเจ็ค**
```bash
git clone https://github.com/Kamonsakgo/Botnoi_Blog.git
cd Botnoi_Blog
```

2. **ติดตั้ง Dependencies**
```bash
go mod tidy
```

3. **ตั้งค่า Environment Variables**
สร้างไฟล์ `.env` และเพิ่มตัวแปรต่อไปนี้:
```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=botnoi_blog
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=your_region
AWS_BUCKET_NAME=your_bucket_name
JWT_SECRET=your_jwt_secret
```

4. **เรียกใช้งานแอปพลิเคชัน**
```bash
go run main.go
```

แอปพลิเคชันจะทำงานที่ `http://localhost:8080`

## API Documentation
เข้าถึง Swagger UI ได้ที่: `http://localhost:8080/api/admin/docs`

## คุณสมบัติหลัก

### 1. การจัดการผู้ใช้ (Users)
- ลงทะเบียนและเข้าสู่ระบบ
- การจัดการโปรไฟล์ผู้ใช้
- ระบบ JWT Authentication

### 2. การจัดการบล็อก (Blogs)
- สร้าง อ่าน แก้ไข และลบบทความ
- อัปโหลดรูปภาพสำหรับบทความ
- การจัดการหมวดหมู่และแท็ก

### 3. กิจกรรมไฮไลท์ (Highlight Events)
- สร้างและจัดการงานกิจกรรม
- อัปโหลดภาพปกกิจกรรม
- การจัดการสถานะกิจกรรม

## การทดสอบ

### รันการทดสอบหน่วย (Unit Tests)
```bash
go test ./src/services/...
```

### รันการทดสอบแบบครบวงจร (Integration Tests)
```bash
go test ./tests/...
```



## การพัฒนา

### โครงสร้าง Clean Architecture
โปรเจ็คนี้ใช้หลักการ Clean Architecture:
- **Entities**: โครงสร้างข้อมูลหลัก
- **Repositories**: Interface สำหรับการเข้าถึงข้อมูล
- **Services**: Business Logic
- **Gateways**: HTTP Handlers

### Middleware
- **JWT Middleware**: ตรวจสอบสิทธิ์การเข้าถึง
- **Logger Middleware**: บันทึกการเรียกใช้ API
- **CORS Middleware**: รองรับการเรียกใช้จาก Domain อื่น

## การมีส่วนร่วม (Contributing)
1. Fork โปรเจ็ค
2. สร้าง Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit การเปลี่ยนแปลง (`git commit -m 'Add some AmazingFeature'`)
4. Push ไปยัง Branch (`git push origin feature/AmazingFeature`)
5. เปิด Pull Request
