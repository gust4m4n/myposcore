-- Dummy FAQ for MyPOS Core
-- Run this after table creation

INSERT INTO faqs (question, answer, category, `order`, is_active, created_at, updated_at) VALUES

-- General Category
('Apa itu MyPOS Core?',
'MyPOS Core adalah sistem Point of Sale (POS) berbasis cloud yang dirancang untuk membantu bisnis mengelola penjualan, inventori, dan operasional toko dengan mudah. Sistem ini mendukung multi-cabang dan multi-pengguna dengan fitur lengkap untuk bisnis retail dan F&B.',
'General',
1,
true,
NOW(),
NOW()
),

('Siapa yang bisa menggunakan MyPOS Core?',
'MyPOS Core cocok untuk berbagai jenis bisnis seperti:
- Toko retail (fashion, elektronik, grocery)
- Restoran dan kafe
- Toko kelontong
- Bisnis dengan multiple cabang
- UMKM yang ingin digitalisasi operasional

Sistem ini dirancang user-friendly sehingga mudah digunakan bahkan untuk yang baru pertama kali menggunakan POS.',
'General',
2,
true,
NOW(),
NOW()
),

('Apakah MyPOS Core tersedia dalam bahasa Indonesia?',
'Ya, MyPOS Core tersedia dalam bahasa Indonesia dan English. Anda dapat mengubah preferensi bahasa di pengaturan akun.',
'General',
3,
true,
NOW(),
NOW()
),

-- Account & Access
('Bagaimana cara membuat akun di MyPOS Core?',
'Untuk membuat akun:
1. Klik tombol "Register" di halaman login
2. Isi informasi bisnis (nama bisnis, alamat, nomor telepon)
3. Buat username dan password
4. Verifikasi email (jika diperlukan)
5. Akun Anda siap digunakan!

Akun pertama otomatis menjadi Super Admin dengan akses penuh.',
'Account',
1,
true,
NOW(),
NOW()
),

('Bagaimana jika lupa password?',
'Jika lupa password:
1. Klik "Forgot Password" di halaman login
2. Masukkan username atau email Anda
3. Ikuti instruksi yang dikirim ke email
4. Buat password baru

Jika masih mengalami kendala, hubungi tim support kami.',
'Account',
2,
true,
NOW(),
NOW()
),

('Berapa banyak user yang bisa ditambahkan?',
'Tidak ada batasan jumlah user. Anda bisa menambahkan sebanyak yang dibutuhkan dengan role berbeda:
- Super Admin (akses penuh)
- Branch Admin (mengelola cabang)
- Cashier (hanya transaksi)
- Staff (terbatas sesuai kebutuhan)',
'Account',
3,
true,
NOW(),
NOW()
),

-- Products & Inventory
('Bagaimana cara menambahkan produk?',
'Cara menambahkan produk:
1. Login ke dashboard
2. Buka menu "Products"
3. Klik tombol "Add Product"
4. Isi detail produk (nama, harga, kategori, stok, dll)
5. Upload foto produk (opsional)
6. Klik "Save"

Anda juga bisa import produk dalam jumlah banyak menggunakan file Excel.',
'Products',
1,
true,
NOW(),
NOW()
),

('Apakah bisa mengelola stok otomatis?',
'Ya, MyPOS Core otomatis mengurangi stok setiap ada transaksi penjualan. Anda juga bisa:
- Set notifikasi low stock
- Tracking history perubahan stok
- Melakukan stock opname
- Transfer stok antar cabang',
'Products',
2,
true,
NOW(),
NOW()
),

('Apakah mendukung barcode?',
'Ya, sistem kami mendukung:
- Scan barcode untuk transaksi cepat
- Generate barcode untuk produk
- Import barcode dari supplier
- Custom barcode format

Anda perlu barcode scanner yang kompatibel (USB atau Bluetooth).',
'Products',
3,
true,
NOW(),
NOW()
),

-- Orders & Payments
('Metode pembayaran apa saja yang didukung?',
'MyPOS Core mendukung berbagai metode pembayaran:
- Cash (Tunai)
- Debit Card
- Credit Card
- E-Wallet (GoPay, OVO, Dana, dll)
- Transfer Bank
- QRIS

Anda bisa mengaktifkan/menonaktifkan metode pembayaran sesuai kebutuhan.',
'Payments',
1,
true,
NOW(),
NOW()
),

('Bagaimana cara melakukan transaksi?',
'Langkah transaksi:
1. Pilih produk dengan scan barcode atau search
2. Tentukan quantity
3. Sistem otomatis kalkulasi total
4. Pilih metode pembayaran
5. Input jumlah bayar
6. Print atau kirim struk digital

Proses sangat cepat, rata-rata hanya 30 detik per transaksi!',
'Orders',
1,
true,
NOW(),
NOW()
),

('Apakah bisa refund transaksi?',
'Ya, refund bisa dilakukan dengan:
1. Akses history transaksi
2. Pilih transaksi yang akan direfund
3. Klik "Refund" dan masukkan alasan
4. Konfirmasi refund
5. Stok otomatis kembali

Refund memerlukan approval dari admin (bisa diatur di settings).',
'Orders',
2,
true,
NOW(),
NOW()
),

-- Reports & Analytics
('Laporan apa saja yang tersedia?',
'MyPOS Core menyediakan berbagai laporan:
- Laporan penjualan (harian, mingguan, bulanan)
- Laporan produk terlaris
- Laporan profit & loss
- Laporan stok
- Laporan kasir performance
- Laporan per cabang

Semua laporan bisa di-export ke Excel atau PDF.',
'Reports',
1,
true,
NOW(),
NOW()
),

('Apakah data laporan real-time?',
'Ya, semua data di dashboard adalah real-time. Setiap transaksi langsung ter-update di:
- Total penjualan hari ini
- Grafik penjualan
- Stok produk
- Revenue

Anda bisa monitoring bisnis dari mana saja, kapan saja.',
'Reports',
2,
true,
NOW(),
NOW()
),

-- Technical
('Apakah perlu internet untuk menggunakan MyPOS Core?',
'Ya, MyPOS Core adalah sistem cloud-based yang memerlukan koneksi internet. Keuntungannya:
- Akses dari mana saja
- Data aman di cloud
- Auto backup
- Update otomatis
- Tidak perlu install software

Rekomendasi minimal: koneksi 5 Mbps untuk operasional lancar.',
'Technical',
1,
true,
NOW(),
NOW()
),

('Perangkat apa saja yang didukung?',
'MyPOS Core bisa diakses melalui:
- Web Browser (Chrome, Firefox, Safari, Edge)
- Windows PC/Laptop
- Mac
- Tablet (Android/iOS)
- Smartphone (untuk monitoring)

Untuk kasir, kami rekomendasikan PC/Laptop atau Tablet untuk pengalaman terbaik.',
'Technical',
2,
true,
NOW(),
NOW()
),

('Apakah data aman?',
'Keamanan data adalah prioritas kami:
- Enkripsi SSL/TLS untuk semua koneksi
- Database terenkripsi
- Backup otomatis setiap hari
- Disaster recovery plan
- Compliance dengan standar keamanan

Data Anda tersimpan di server terpercaya dengan uptime 99.9%.',
'Technical',
3,
true,
NOW(),
NOW()
),

-- Pricing & Subscription
('Berapa biaya berlangganan MyPOS Core?',
'Kami menawarkan berbagai paket:
- Starter: Untuk bisnis kecil (1 cabang)
- Business: Untuk bisnis menengah (hingga 5 cabang)
- Enterprise: Untuk bisnis besar (unlimited cabang)

Hubungi tim sales kami untuk penawaran khusus dan trial gratis 14 hari!',
'Pricing',
1,
true,
NOW(),
NOW()
),

('Apakah ada biaya setup atau training?',
'Tidak ada biaya setup! Kami menyediakan:
- Free onboarding session
- Video tutorial lengkap
- Dokumentasi detail
- Customer support 24/7

Training tambahan bisa diatur sesuai kebutuhan (opsional).',
'Pricing',
2,
true,
NOW(),
NOW()
),

-- Support
('Bagaimana cara menghubungi support?',
'Anda bisa menghubungi kami melalui:
- Email: support@myposcore.com
- WhatsApp: +62 812-3456-7890
- Live Chat di dashboard
- Help Center (dokumentasi online)

Tim support kami siap membantu 24/7!',
'Support',
1,
true,
NOW(),
NOW()
),

('Apakah ada dokumentasi atau tutorial?',
'Ya, kami menyediakan:
- Video tutorial untuk setiap fitur
- User manual lengkap (PDF)
- Knowledge base online
- Webinar bulanan
- Blog dengan tips & tricks

Semua bisa diakses gratis di Help Center kami.',
'Support',
2,
true,
NOW(),
NOW()
);
