-- Master Seed Data
-- Case Study: Kabupaten Subang

INSERT INTO districts
(id, name)
VALUES
(1, 'Subang'),
(2, 'Pamanukan'),
(3, 'Ciasem'),
(4, 'Pusakanagara'),
(5, 'Pusakajaya'),
(6, 'Blanakan'),
(7, 'Legonkulon'),
(8, 'Compreng'),
(9, 'Binong'),
(10, 'Tambakdahan'),
(11, 'Pagaden'),
(12, 'Pagaden Barat'),
(13, 'Cipunagara'),
(14, 'Purwadadi'),
(15, 'Cikaum'),
(16, 'Patokbeusi'),
(17, 'Kalijati'),
(18, 'Dawuan'),
(19, 'Cipeundeuy'),
(20, 'Pabuaran'),
(21, 'Cibogo'),
(22, 'Cijambe'),
(23, 'Kasomalang'),
(24, 'Jalancagak'),
(25, 'Ciater'),
(26, 'Sagalaherang'),
(27, 'Serangpanjang'),
(28, 'Cisalak'),
(29, 'Tanjungsiang');

INSERT INTO departments
(id, name, code)
VALUES
(1,'Sekretariat Daerah','SETDA'),
(2,'Sekretariat DPRD','SETWAN'),
(3,'Inspektorat','INSPEKTORAT'),
(4,'Badan Perencanaan Pembangunan Penelitian Dan Pengembangan Daerah','BP4D'),
(5,'Badan Kepegawaian dan Pengembangan Sumber Daya Manusia','BKPSDM'),
(6,'Badan Keuangan dan Aset Daerah','BKAD'),
(7,'Badan Pendapatan Daerah','BAPENDA'),
(8,'Dinas Komunikasi dan Informatika','DISKOMINFO'),
(9,'Dinas Pendidikan','DISDIK'),
(10,'Dinas Kesehatan','DINKES'),
(11,'Dinas Sosial','DINSOS'),
(12,'Dinas Kependudukan dan Pencatatan Sipil','DISDUKCAPIL'),
(13,'Dinas Pemberdayaan Masyarakat dan Desa','DPMD'),
(14,'Dinas Pengendalian Penduduk, KB, Pemberdayaan Perempuan dan Perlindungan Anak','DP2KBP3A'),
(15,'Dinas Tenaga Kerja','DISNAKER'),
(16,'Dinas Penanaman Modal Dan Pelayanan Terpadu Satu Pintu','DPMPTSP'),
(17,'Dinas Perdagangan dan Perindustrian','DISDAGIN'),
(18,'Dinas Koperasi UKM','DKUKM'),
(19,'Dinas Ketahanan Pangan','DKP'),
(20,'Dinas Pertanian','DISTAN'),
(21,'Dinas Perikanan','DISKAN'),
(22,'Dinas Peternakan','DISNAK'),
(23,'Dinas Lingkungan Hidup','DLH'),
(24,'Dinas Pekerjaan Umum dan Penataan Ruang','PUPR'),
(25,'Dinas Perhubungan','DISHUB'),
(26,'Dinas Perumahan, Kawasan Permukiman dan Pertanahan','DISPERKIMTAN'),
(27,'Pemadam Kebakaran','DAMKAR'),
(28,'Badan Penanggulangan Bencana Daerah','BPBD'),
(29,'Satpol PP','SATPOLPP'),
(30,'Dinas Pariwisata, Kepemudaan dan Olah Raga','DISPARPORA');

INSERT INTO categories
(id, department_id, name, description)
VALUES
(1, 24, 'Jalan Rusak', 'Kerusakan jalan umum yang mengganggu mobilitas masyarakat'),
(2, 24, 'Jembatan Rusak', 'Kerusakan jembatan yang membahayakan pengguna jalan'),
(3, 24, 'Drainase Tersumbat', 'Saluran drainase tersumbat yang berpotensi menyebabkan banjir'),
(4, 24, 'Trotoar Rusak', 'Trotoar atau jalur pejalan kaki mengalami kerusakan'),

(5, 25, 'Lampu Jalan Mati', 'Lampu penerangan jalan umum tidak berfungsi'),
(6, 25, 'Marka Jalan Rusak', 'Marka jalan memudar atau rusak'),
(7, 25, 'Rambu Lalu Lintas Rusak', 'Rambu lalu lintas hilang atau rusak'),
(8, 25, 'Kemacetan', 'Kemacetan lalu lintas yang membutuhkan penanganan Dishub'),

(9, 23, 'Sampah Menumpuk', 'Tumpukan sampah yang mengganggu kebersihan lingkungan'),
(10, 23, 'Pohon Tumbang', 'Pohon tumbang atau pohon yang berpotensi membahayakan'),
(11, 23, 'Taman Kota Rusak', 'Kerusakan taman kota atau ruang terbuka hijau'),

(12, 28, 'Banjir', 'Genangan atau banjir yang mengganggu aktivitas masyarakat'),
(13, 28, 'Longsor', 'Bencana tanah longsor'),
(14, 28, 'Bangunan Roboh', 'Bangunan roboh akibat bencana atau faktor lainnya'),

(15, 27, 'Kebakaran', 'Kejadian kebakaran pada bangunan atau lahan'),

(16, 20, 'Hewan Terlantar', 'Hewan liar atau ternak yang mengganggu lingkungan'),

(17, 29, 'PKL Mengganggu Ketertiban', 'Pedagang kaki lima yang melanggar ketertiban umum'),
(18, 29, 'Vandalisme Fasilitas Umum', 'Perusakan fasilitas umum oleh oknum'),

(19, 24, 'Air Bersih', 'Gangguan distribusi atau kualitas air bersih'),

(20, 23, 'Limbah Liar', 'Pembuangan limbah secara ilegal yang mencemari lingkungan');

INSERT INTO emergency_contacts
(id, department_id, district_id, name, phone_number, description)
VALUES
(1, NULL, NULL, 'Call Center Darurat', '112', 'Layanan panggilan darurat nasional'),
(2, 27, NULL, 'Pemadam Kebakaran Subang', '(0260) 411007', 'Layanan pemadam kebakaran'),
(3, 28, NULL, 'BPBD Kabupaten Subang', '081224334343', 'Pelaporan bencana alam'),
(4, 10, NULL, 'Ambulans PSC 119', '082152119119', 'Pelayanan ambulans darurat'),
(5, NULL, NULL, 'Polres Subang', '08131550110', 'Layanan kepolisian'),
(6, 29, NULL, 'Satpol PP Kabupaten Subang', '(0260) 411234', 'Ketertiban umum'),
(7, NULL, NULL, 'PLN Pengaduan', '(0260) 421581', 'Gangguan listrik'),
(8, NULL, NULL, 'PDAM Tirta Rangga', '081908197024', 'Gangguan air bersih'),
(9, 10, NULL, 'RSUD Subang', '(0260) 411421', 'Rumah sakit umum daerah'),
(10, NULL, NULL, 'PMI Kabupaten Subang', '(0260) 411985', 'Palang Merah Indonesia');

INSERT INTO regional_news
(id, district_id, title, content, category, scope)
VALUES

(
1,
NULL,
'Program Perbaikan Jalan Kabupaten Dimulai',
'Pemerintah Kabupaten Subang memulai program perbaikan jalan kabupaten secara bertahap untuk meningkatkan keselamatan dan kenyamanan pengguna jalan.',
'announcement',
'global'
),

(
2,
NULL,
'BPBD Mengimbau Warga Waspada Cuaca Ekstrem',
'BPBD Kabupaten Subang mengimbau masyarakat agar meningkatkan kewaspadaan terhadap potensi hujan lebat disertai angin kencang selama beberapa hari ke depan.',
'emergency',
'global'
),

(
3,
1,
'Normalisasi Drainase Kecamatan Subang',
'Pekerjaan normalisasi drainase dilakukan untuk mengurangi potensi genangan saat musim hujan.',
'news',
'district'
),

(
4,
2,
'Perbaikan Lampu Jalan Jalur Pantura Pamanukan',
'Perbaikan lampu penerangan jalan umum dilakukan secara bertahap pada ruas jalan nasional wilayah Pamanukan.',
'announcement',
'district'
),

(
5,
25,
'Objek Wisata Ciater Tetap Beroperasi',
'Pengunjung diimbau menjaga kebersihan kawasan wisata dan mematuhi peraturan yang berlaku.',
'news',
'district'
),

(
6,
NULL,
'Festival Nanas Subang 2026 Akan Digelar Bulan Depan',
'Festival tahunan yang menampilkan UMKM, pertanian, dan budaya lokal akan kembali diselenggarakan.',
'event',
'global'
),

(
7,
NULL,
'Layanan Administrasi Tetap Beroperasi Saat Libur Nasional',
'Pelayanan administrasi kependudukan tertentu tetap dibuka melalui layanan digital.',
'announcement',
'global'
),

(
8,
NULL,
'Pemangkasan Pohon Rawan Tumbang Dilaksanakan Bertahap',
'Dinas Lingkungan Hidup melakukan pemangkasan pohon yang berpotensi membahayakan pengguna jalan.',
'news',
'global'
);