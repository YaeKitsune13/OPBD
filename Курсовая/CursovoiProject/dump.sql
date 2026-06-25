/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-12.3.2-MariaDB, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: vet_clinic
-- ------------------------------------------------------
-- Server version	8.0.46

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `appointment_diagnoses`
--

DROP TABLE IF EXISTS `appointment_diagnoses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointment_diagnoses` (
  `appointment_id` bigint unsigned NOT NULL,
  `diagnosis_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`appointment_id`,`diagnosis_id`),
  KEY `fk_appointment_diagnoses_diagnosis` (`diagnosis_id`),
  CONSTRAINT `fk_appointment_diagnoses_appointment` FOREIGN KEY (`appointment_id`) REFERENCES `appointments` (`id`),
  CONSTRAINT `fk_appointment_diagnoses_diagnosis` FOREIGN KEY (`diagnosis_id`) REFERENCES `diagnoses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointment_diagnoses`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `appointment_diagnoses` WRITE;
/*!40000 ALTER TABLE `appointment_diagnoses` DISABLE KEYS */;
INSERT INTO `appointment_diagnoses` VALUES
(11,1),
(13,2),
(13,4),
(10,7),
(10,9),
(10,10),
(10,11);
/*!40000 ALTER TABLE `appointment_diagnoses` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `appointments`
--

DROP TABLE IF EXISTS `appointments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `client_id` bigint unsigned DEFAULT NULL,
  `pet_id` bigint unsigned DEFAULT NULL,
  `doctor_id` bigint unsigned DEFAULT NULL,
  `service_id` bigint unsigned DEFAULT NULL,
  `scheduled_at` datetime(3) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'waiting',
  `comment` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_appointments_deleted_at` (`deleted_at`),
  KEY `fk_appointments_doctor` (`doctor_id`),
  KEY `fk_appointments_service` (`service_id`),
  KEY `fk_pets_appointments` (`pet_id`),
  KEY `fk_appointments_client` (`client_id`),
  CONSTRAINT `fk_appointments_client` FOREIGN KEY (`client_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_appointments_doctor` FOREIGN KEY (`doctor_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_appointments_service` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`),
  CONSTRAINT `fk_pets_appointments` FOREIGN KEY (`pet_id`) REFERENCES `pets` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointments`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `appointments` WRITE;
/*!40000 ALTER TABLE `appointments` DISABLE KEYS */;
INSERT INTO `appointments` VALUES
(1,'2026-05-24 16:34:47.977','2026-05-24 17:09:17.115',NULL,2,4,3,1,'2026-05-24 09:00:00.000','done',''),
(2,'2026-05-24 16:35:38.532','2026-05-24 17:09:19.765',NULL,2,4,3,1,'2026-05-24 10:00:00.000','rejected',''),
(3,'2026-05-24 16:36:38.091','2026-05-25 21:15:44.155',NULL,2,4,3,2,'2026-05-24 10:00:00.000','done',''),
(4,'2026-05-25 16:20:13.848','2026-05-25 23:48:29.741',NULL,1,1,4,1,'2026-05-25 09:00:00.000','done','Подозрительно ведёт себя, грызётся'),
(5,'2026-05-25 16:30:44.568','2026-05-25 21:14:47.357',NULL,1,1,3,1,'2026-05-25 09:00:00.000','rejected',''),
(6,'2026-05-25 16:36:59.915','2026-05-25 23:48:47.073',NULL,1,1,4,1,'2026-05-25 10:00:00.000','done','А'),
(7,'2026-05-25 16:39:39.147','2026-05-25 23:54:48.130',NULL,1,3,4,1,'2026-05-25 11:00:00.000','done',''),
(8,'2026-05-25 23:38:12.035','2026-05-25 23:42:16.726',NULL,1,1,3,1,'2026-05-26 09:00:00.000','done',''),
(9,'2026-05-26 10:21:44.075','2026-05-26 10:23:12.214',NULL,1,1,3,1,'2026-05-26 10:00:00.000','done',''),
(10,'2026-06-19 23:21:01.211','2026-06-19 23:23:56.051',NULL,1,1,3,1,'2026-06-19 09:00:00.000','done',''),
(11,'2026-06-19 23:31:06.040','2026-06-19 23:31:59.900',NULL,1,1,3,1,'2026-06-19 10:00:00.000','done',''),
(12,'2026-06-19 23:38:34.288','2026-06-19 23:38:56.939',NULL,1,1,3,1,'2026-06-19 11:00:00.000','confirmed',''),
(13,'2026-06-21 13:06:44.117','2026-06-21 13:17:48.755',NULL,7,6,3,1,'2026-06-22 10:00:00.000','done',''),
(14,'2026-06-21 13:13:21.231','2026-06-21 13:13:32.922',NULL,7,6,3,1,'2026-06-22 11:00:00.000','rejected','');
/*!40000 ALTER TABLE `appointments` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `cart_items`
--

DROP TABLE IF EXISTS `cart_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `cart_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `quantity` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cart_items_deleted_at` (`deleted_at`),
  KEY `fk_cart_items_product` (`product_id`),
  CONSTRAINT `fk_cart_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cart_items`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `cart_items` WRITE;
/*!40000 ALTER TABLE `cart_items` DISABLE KEYS */;
INSERT INTO `cart_items` VALUES
(1,'2026-05-24 17:09:57.772','2026-05-24 17:09:57.772','2026-05-24 17:10:31.022',2,1,1),
(2,'2026-05-24 17:10:02.589','2026-05-24 17:10:06.258','2026-05-24 17:10:06.434',2,8,1),
(3,'2026-05-24 17:10:09.902','2026-05-24 17:10:09.902','2026-05-24 17:10:31.022',2,3,1),
(4,'2026-05-25 16:20:38.755','2026-05-25 16:20:38.755','2026-05-25 16:20:48.659',1,1,1),
(5,'2026-05-25 16:20:39.402','2026-05-25 16:20:39.402','2026-05-25 16:20:48.659',1,7,1),
(6,'2026-05-25 16:20:40.005','2026-05-25 16:20:40.005','2026-05-25 16:20:48.659',1,3,1),
(7,'2026-05-25 21:20:44.063','2026-05-25 21:20:44.202','2026-05-25 21:20:46.461',1,2,2),
(8,'2026-05-25 21:20:56.200','2026-05-25 21:20:56.200','2026-05-25 21:20:59.351',1,3,1),
(9,'2026-05-25 21:20:56.606','2026-05-25 21:20:56.606','2026-05-25 21:20:59.351',1,2,1),
(10,'2026-05-25 21:20:56.965','2026-05-25 21:20:56.965','2026-05-25 21:20:59.351',1,1,1),
(11,'2026-05-25 21:54:58.824','2026-05-25 21:54:58.824','2026-05-25 21:55:12.551',1,4,1),
(12,'2026-05-25 21:55:03.905','2026-05-25 21:55:11.373','2026-05-25 21:55:24.101',1,1,12),
(13,'2026-05-25 21:55:04.782','2026-05-25 21:55:04.782','2026-05-25 21:55:15.773',1,2,1),
(14,'2026-05-25 23:45:51.714','2026-05-25 23:46:05.702','2026-05-25 23:46:18.444',1,1,1),
(15,'2026-05-25 23:45:55.041','2026-05-25 23:45:58.164','2026-05-25 23:46:18.444',1,4,2),
(16,'2026-05-25 23:45:55.583','2026-05-25 23:45:58.643','2026-05-25 23:46:18.444',1,8,2),
(17,'2026-05-25 23:45:56.019','2026-05-25 23:45:56.019','2026-05-25 23:46:18.444',1,7,1);
/*!40000 ALTER TABLE `cart_items` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `diagnoses`
--

DROP TABLE IF EXISTS `diagnoses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `diagnoses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_diagnoses_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `diagnoses`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `diagnoses` WRITE;
/*!40000 ALTER TABLE `diagnoses` DISABLE KEYS */;
INSERT INTO `diagnoses` VALUES
(1,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Клинически здоров'),
(2,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Острый гастроэнтерит'),
(3,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Мочекаменная болезнь (МКБ)'),
(4,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Хроническая почечная недостаточность (ХПН)'),
(5,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Отит наружного уха'),
(6,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Пироплазмоз (бабезиоз)'),
(7,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Блошиный аллергический дерматит'),
(8,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Парвовирусный энтерит'),
(9,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Панлейкопения кошек'),
(10,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Инородное тело в ЖКТ'),
(11,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Разрыв передней крестообразной связки'),
(12,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Зубной камень'),
(13,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Малокклюзия'),
(14,'2026-06-19 18:23:14.000','2026-06-19 18:23:14.000',NULL,'Калицивироз');
/*!40000 ALTER TABLE `diagnoses` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `medical_protocols`
--

DROP TABLE IF EXISTS `medical_protocols`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `medical_protocols` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `appointment_id` bigint unsigned DEFAULT NULL,
  `weight_at_visit` double DEFAULT NULL,
  `treatment` longtext,
  `medications` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_medical_protocols_appointment_id` (`appointment_id`),
  KEY `idx_medical_protocols_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_appointments_protocol` FOREIGN KEY (`appointment_id`) REFERENCES `appointments` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `medical_protocols`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `medical_protocols` WRITE;
/*!40000 ALTER TABLE `medical_protocols` DISABLE KEYS */;
INSERT INTO `medical_protocols` VALUES
(1,'2026-05-24 17:09:17.114','2026-05-24 17:09:17.114',NULL,1,4.5,'Диета, покой','Омепразол'),
(2,'2026-05-25 21:15:44.154','2026-05-25 21:15:44.154',NULL,3,4.6,'Мыть чаще и давать таблетки от клещей','Таблетки от клещей'),
(3,'2026-05-25 23:42:16.726','2026-05-25 23:42:16.726',NULL,8,4.2,'Помощь в линьке ежедневный уход','Таблитки от линьки'),
(4,'2026-05-25 23:48:29.741','2026-05-25 23:48:29.741',NULL,4,4.7,'ухаживать','таблетки'),
(5,'2026-05-25 23:48:47.072','2026-05-25 23:48:47.072',NULL,6,8,'Рекомендаций\n','Лекартсва'),
(6,'2026-05-25 23:54:48.129','2026-05-25 23:54:48.129',NULL,7,4.8,'Регулярная мойка и таблеки','Таблетки для кошек'),
(7,'2026-05-26 10:23:12.213','2026-05-26 10:23:12.213',NULL,9,4.9,'Очиска кишечника','Таблетки от линьки'),
(8,'2026-06-19 23:23:56.048','2026-06-19 23:23:56.048',NULL,10,4.5,'Уход','Лекарство 1мл 5д'),
(9,'2026-06-19 23:31:59.897','2026-06-19 23:31:59.897',NULL,11,3.5,'Чаще мыть','Лекарство 1мл 5д'),
(10,'2026-06-21 13:17:48.751','2026-06-21 13:17:48.751',NULL,13,5.4,'Уход','Лекарства');
/*!40000 ALTER TABLE `medical_protocols` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `order_items`
--

DROP TABLE IF EXISTS `order_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `quantity` bigint unsigned DEFAULT NULL,
  `price` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_items_deleted_at` (`deleted_at`),
  KEY `fk_order_items_product` (`product_id`),
  KEY `fk_orders_items` (`order_id`),
  CONSTRAINT `fk_order_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_orders_items` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_items`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `order_items` WRITE;
/*!40000 ALTER TABLE `order_items` DISABLE KEYS */;
INSERT INTO `order_items` VALUES
(1,'2026-05-24 17:10:31.019','2026-05-24 17:10:31.019',NULL,1,1,1,450),
(2,'2026-05-24 17:10:31.019','2026-05-24 17:10:31.019',NULL,1,3,1,3200),
(3,'2026-05-25 16:20:48.656','2026-05-25 16:20:48.656',NULL,2,1,1,450),
(4,'2026-05-25 16:20:48.656','2026-05-25 16:20:48.656',NULL,2,7,1,650),
(5,'2026-05-25 16:20:48.656','2026-05-25 16:20:48.656',NULL,2,3,1,3200),
(6,'2026-05-25 21:20:46.459','2026-05-25 21:20:46.459',NULL,3,2,2,980),
(7,'2026-05-25 21:20:59.349','2026-05-25 21:20:59.349',NULL,4,3,1,3200),
(8,'2026-05-25 21:20:59.349','2026-05-25 21:20:59.349',NULL,4,2,1,980),
(9,'2026-05-25 21:20:59.349','2026-05-25 21:20:59.349',NULL,4,1,1,450),
(10,'2026-05-25 23:46:18.442','2026-05-25 23:46:18.442',NULL,5,1,1,450),
(11,'2026-05-25 23:46:18.442','2026-05-25 23:46:18.442',NULL,5,4,2,720),
(12,'2026-05-25 23:46:18.442','2026-05-25 23:46:18.442',NULL,5,8,2,380),
(13,'2026-05-25 23:46:18.442','2026-05-25 23:46:18.442',NULL,5,7,1,650);
/*!40000 ALTER TABLE `order_items` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `total_amount` bigint unsigned DEFAULT NULL,
  `status` enum('paid','confirmed','delivered') DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES
(1,'2026-05-24 17:10:31.018','2026-05-24 17:10:31.018',NULL,2,3650,'paid'),
(2,'2026-05-25 16:20:48.656','2026-05-25 16:20:48.656',NULL,1,4300,'paid'),
(3,'2026-05-25 21:20:46.459','2026-05-25 21:20:46.459',NULL,1,1960,'paid'),
(4,'2026-05-25 21:20:59.349','2026-05-25 21:20:59.349',NULL,1,4630,'paid'),
(5,'2026-05-25 23:46:18.441','2026-05-25 23:46:18.441',NULL,1,3300,'paid');
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `pets`
--

DROP TABLE IF EXISTS `pets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `owner_id` bigint unsigned DEFAULT NULL,
  `name` longtext,
  `species` longtext,
  `breed` longtext,
  `birth_date` datetime(3) DEFAULT NULL,
  `weight` double DEFAULT NULL,
  `avatar` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_pets_deleted_at` (`deleted_at`),
  KEY `fk_users_pets` (`owner_id`),
  CONSTRAINT `fk_users_pets` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pets`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `pets` WRITE;
/*!40000 ALTER TABLE `pets` DISABLE KEYS */;
INSERT INTO `pets` VALUES
(1,'2026-05-24 16:26:33.224','2026-05-25 16:11:08.011',NULL,1,'Барсик','Кошка','Британская','2020-03-15 05:00:00.000',4.5,''),
(2,'2026-05-24 16:26:55.258','2026-05-24 16:26:55.258','2026-05-24 16:27:21.217',1,'Рекс','Собака','Лабрадор','2019-07-22 05:00:00.000',28,''),
(3,'2026-05-24 16:27:19.537','2026-05-24 16:27:19.537',NULL,1,'Мурка','Кошка','Персидская','2021-01-10 05:00:00.000',3.8,''),
(4,'2026-05-24 16:27:52.639','2026-05-25 21:16:56.759',NULL,2,'Рекс','Собака','Лабрадор','2019-07-22 05:00:00.000',28,''),
(5,'2026-05-25 23:37:13.219','2026-05-25 23:37:13.219','2026-05-25 23:37:39.772',1,'dasd','Кошка','asdas','2026-05-13 05:00:00.000',2,''),
(6,'2026-06-21 13:03:09.866','2026-06-21 13:03:09.866',NULL,7,'Барсик','Кошка','Британская','2021-05-12 05:00:00.000',4.2,'');
/*!40000 ALTER TABLE `pets` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `description` longtext,
  `price` bigint DEFAULT NULL,
  `category` longtext,
  `stock` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES
(1,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Амоксициллин 15%','Антибиотик широкого спектра действия, суспензия для инъекций 10мл',450,'Антибиотики',50),
(2,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Фронтлайн Комбо','Капли на холку от блох и клещей для кошек',980,'Паразиты',100),
(3,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Бравекто (10-20 кг)','Таблетка от клещей и блох для собак, защита на 12 недель',3200,'Паразиты',30),
(4,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Дронтал для кошек','Антигельминтик (от глистов), 1 таблетка на 4 кг веса',720,'Паразиты',150),
(5,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Мелоксикам 0.2%','Обезболивающее и противовоспалительное средство для кошек и собак',560,'Обезболивающие',40),
(6,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Витамины 8in1 Excel','Мультивитаминный комплекс для взрослых собак, 70 таб.',1100,'Витамины',25),
(7,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Беафар Лавиета','Витамины для улучшения шерсти у кошек, капли 50мл',650,'Витамины',35),
(8,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Котэрвин','Средство для профилактики и лечения мочекаменной болезни кошек',380,'Терапия',60),
(9,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Стоп-зуд суспензия','Для лечения воспалительных и аллергических заболеваний кожи',490,'Дерматология',45),
(10,'2026-05-21 11:27:26.000','2026-05-21 11:27:26.000',NULL,'Хлоргексидин 0.05%','Антисептический раствор для обработки ран, 100мл',80,'Антисептики',200);
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `services`
--

DROP TABLE IF EXISTS `services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `services` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `price` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_services_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `services`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `services` WRITE;
/*!40000 ALTER TABLE `services` DISABLE KEYS */;
INSERT INTO `services` VALUES
(1,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Первичный полный осмотр и консультация',1200),
(2,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Профилактический осмотр перед вакцинацией',600),
(3,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Стрижка когтей (кошки и мелкие собаки)',450),
(4,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Гигиеническая чистка ушных раковин',400),
(5,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Обработка глаз и очистка слезных дорожек',350),
(6,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Вычесывание подшерстка (экспресс-линька)',1800),
(7,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Удаление колтунов и гигиеническая стрижка',1500),
(8,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Санация ротовой полости (чистка зубов пастой)',800),
(9,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Мытье питомца и сушка феном',1200),
(10,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Обработка от кожных паразитов (нанесение капель)',250),
(11,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Чипирование с внесением в базу данных',1600),
(12,'2026-05-21 11:27:16.000','2026-05-21 11:27:16.000',NULL,'Консультация по питанию и воспитанию кошки',900);
/*!40000 ALTER TABLE `services` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `password` longtext,
  `first_name` longtext,
  `last_name` longtext,
  `phone` varchar(191) DEFAULT NULL,
  `role` enum('client','doctor') DEFAULT NULL,
  `specialization` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  UNIQUE KEY `uni_users_phone` (`phone`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES
(1,'2026-05-24 16:23:47.672','2026-06-21 13:42:34.778',NULL,'ivanova@gmail.com','$2a$12$IPSJiMkQ19/jPvtk..tQd.8ffNcCPOgJeazTwe9AFTYvVGqs/.3zO','Мария','Иванова','+79161234567','client',''),
(2,'2026-05-24 16:24:11.219','2026-05-24 16:24:11.219',NULL,'petrov@gmail.com','$2a$12$kQA1V04lajDCIDQcJUGLAuv7L55cAt8WHBd1R/pvdCtDAxJHwcE8q','Алексей','Петров','+79261234568','client',''),
(3,'2026-05-24 16:24:41.454','2026-05-24 16:24:41.454',NULL,'doktor1@vetclinic.ru','$2a$12$JzlSN1le0U.8QSa0PI78KesSYTMaUNNQRQdPWYfYzTh2k0WeKtN4a','Иван','Смирнов','+79301234569','doctor','Терапевт'),
(4,'2026-05-24 16:25:00.237','2026-05-24 16:25:00.237',NULL,'doktor2@vetclinic.ru','$2a$12$4eyZh/9wzR3udfbHs2JyIOmLPGLeKzYMcqWVj/kzLOcxlWS9li0KW','Ольга','Козлова','+79401234570','doctor','Хирург'),
(5,'2026-05-25 16:17:21.810','2026-05-25 16:17:21.810',NULL,'aep@gmail.com','$2a$12$G4s0LJLbpfCvVVvJDVFCKurZHgVinqZDUXsav0cKIKAuxsDyV7HXa','Никита','Петрович','+79001235939','client',''),
(7,'2026-06-21 12:53:55.955','2026-06-21 12:53:55.955',NULL,'ivanov@gmail.com','$2a$12$TMRSXB9jj5To/zHoqvGDnOJ5e/MlL692MUvkb9/euBcuNZ3QUVRmm','Иван','Иванов','+79991112233','client','');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2026-06-25 23:39:55
