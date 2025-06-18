

# API สำหรับจัดการบล็อกและเหตุการณ์สำคัญ (Highlight Event) ด้วย Fiber

โปรเจกต์นี้เป็น RESTful API ที่พัฒนาโดยใช้ Go (Golang) และ Fiber Framework ซึ่งใช้ในการจัดการบล็อกและเหตุการณ์สำคัญ โดยสามารถสร้าง อ่าน แก้ไข และลบบล็อกและเหตุการณ์สำคัญ รวมถึงการจัดการการอัพโหลดรูปภาพไปยัง S3 และการตรวจสอบสิทธิ์การเข้าถึงข้อมูลต่างๆ ด้วย JWT

## ฟีเจอร์

- **การจัดการผู้ใช้**: สร้าง, อ่าน, และจัดการผู้ใช้โดยใช้สิทธิ์ที่แตกต่างกัน
- **การจัดการบล็อก**: สร้าง, แก้ไข, ลบ และดึงข้อมูลบล็อก รวมถึงการอัพโหลดรูปภาพ
- **การจัดการเหตุการณ์สำคัญ**: การจัดการเหตุการณ์สำคัญโดยสามารถเพิ่ม, แก้ไข, ลบ และดึงข้อมูลเหตุการณ์สำคัญ รวมถึงการอัพโหลดรูปภาพ
- **การยืนยันตัวตนด้วย JWT**: ป้องกันการเข้าถึงบางเส้นทางโดยต้องมี JWT ในการยืนยันตัวตน

## การติดตั้ง

### สิ่งที่ต้องการ

- Go 1.18 หรือสูงกว่า
- MongoDB
- S3 Bucket (สำหรับการอัพโหลดรูปภาพ)

### ขั้นตอนการติดตั้ง

1. โคลนโปรเจกต์:
   ```bash
   git clone https://github.com/Kamonsakgo/Botnoi_Blog.git
   cd Botnoi_Blog
   ```

2. ติดตั้ง dependencies:
   ```bash
   go mod tidy
   ```

3. สร้างไฟล์ `.env` ในโฟลเดอร์หลักและกำหนดค่าตัวแปรต่างๆ ดังนี้:
   ```env
   PORT=8080
   MONGO_URI=mongodb://localhost:27017
   S3_BUCKET_NAME=your-bucket-name
   S3_ACCESS_KEY=your-access-key
   S3_SECRET_KEY=your-secret-key
   ```

4. โหลดตัวแปรจากไฟล์ `.env`:
   ควรโหลดตัวแปรเหล่านี้โดยใช้ `godotenv` เพื่อให้ตัวแปรสามารถใช้งานได้ในโปรเจกต์

5. รันแอปพลิเคชัน:
   ```bash
   go run main.go
   ```

   แอปพลิเคชันจะเริ่มทำงานที่พอร์ตที่กำหนดในตัวแปร `PORT`

## รายการ API Endpoint

### การจัดการผู้ใช้

- `POST /api/v1/users/add_user`: สร้างผู้ใช้ใหม่
- `GET /api/v1/users/users`: ดึงข้อมูลผู้ใช้ทั้งหมด

### การจัดการบล็อก

- `GET /api/blog/get_all_blog`: ดึงบล็อกทั้งหมด
- `GET /api/blog/get_blog`: ดึงบล็อกเฉพาะตาม ID
- `POST /api/blog/insert_blog`: สร้างบล็อกใหม่ (ต้องใช้ JWT)
- `PUT /api/blog/update_blog`: แก้ไขบล็อกที่มีอยู่ (ต้องใช้ JWT)
- `DELETE /api/blog/delete_blog`: ลบบล็อกตาม ID (ต้องใช้ JWT)
- `POST /api/blog/upload_image`: อัพโหลดรูปภาพให้กับบล็อก (ต้องใช้ JWT)

### การจัดการเหตุการณ์สำคัญ (Highlight Event)

- `GET /api/highlightevent/get_all`: ดึงเหตุการณ์สำคัญทั้งหมด
- `GET /api/highlightevent/get_one`: ดึงเหตุการณ์สำคัญตาม ID
- `POST /api/highlightevent/insert`: เพิ่มเหตุการณ์สำคัญใหม่ (ต้องใช้ JWT)
- `PUT /api/highlightevent/update`: แก้ไขเหตุการณ์สำคัญที่มีอยู่ (ต้องใช้ JWT)
- `DELETE /api/highlightevent/delete`: ลบเหตุการณ์สำคัญตาม ID (ต้องใช้ JWT)

## การยืนยันตัวตนด้วย JWT

API ใช้ JWT ในการยืนยันตัวตนเพื่อเข้าถึงบางเส้นทาง (เช่น การสร้าง, แก้ไข, หรือลบข้อมูล) โดยให้ส่ง token ใน header `Authorization` ในรูปแบบ:

```
Authorization: Bearer <your-jwt-token>
```

## การจัดการข้อผิดพลาด

ทุกข้อผิดพลาดจะส่งกลับในรูปแบบ `JSON` โดยมีโครงสร้างดังนี้:

```json
{
  "message": "ข้อความข้อผิดพลาด"
}
```

## ไลเซนส์

โปรเจกต์นี้ใช้ MIT License - ดูรายละเอียดในไฟล์ [LICENSE](LICENSE)

## ขอบคุณ

- [Fiber Framework](https://github.com/gofiber/fiber) - เว็บเฟรมเวิร์คสำหรับ Go
- [MongoDB](https://www.mongodb.com) - ฐานข้อมูล NoSQL
- [S3](https://aws.amazon.com/s3/) - บริการคลาวด์สำหรับจัดเก็บรูปภาพ
