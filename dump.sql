-- MySQL dump 10.13  Distrib 8.0.37, for Linux (x86_64)
--
-- Host: localhost    Database: Library
-- ------------------------------------------------------
-- Server version	8.0.37-0ubuntu0.22.04.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `BookRequests`
--

DROP TABLE IF EXISTS `BookRequests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `BookRequests` (
  `RequestID` int NOT NULL AUTO_INCREMENT,
  `UserID` int DEFAULT NULL,
  `BookID` int DEFAULT NULL,
  `RequestDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Status` enum('Pending','Accepted','Returned','Denied') DEFAULT 'Pending',
  `AcceptDate` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`RequestID`),
  KEY `UserID` (`UserID`),
  KEY `BookRequests_ibfk_2` (`BookID`),
  CONSTRAINT `BookRequests_ibfk_1` FOREIGN KEY (`UserID`) REFERENCES `User` (`UserID`),
  CONSTRAINT `BookRequests_ibfk_2` FOREIGN KEY (`BookID`) REFERENCES `Books` (`BookID`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `BookRequests`
--

LOCK TABLES `BookRequests` WRITE;
/*!40000 ALTER TABLE `BookRequests` DISABLE KEYS */;
INSERT INTO `BookRequests` VALUES (44,28,17,'2024-06-07 14:29:44','Returned','2024-06-07 14:30:15'),(48,30,17,'2024-06-08 09:00:40','Denied',NULL),(54,31,17,'2024-06-08 22:41:52','Denied',NULL),(55,31,22,'2024-06-08 22:41:52','Returned','2024-06-08 22:50:05'),(56,31,23,'2024-06-08 22:41:52','Denied',NULL),(57,31,17,'2024-06-08 22:53:18','Pending',NULL),(58,31,22,'2024-06-08 22:53:18','Pending',NULL),(59,31,23,'2024-06-08 22:53:18','Denied',NULL),(60,41,24,'2024-07-03 12:23:41','Returned','2024-07-05 21:58:11'),(61,41,22,'2024-07-06 00:33:02','Pending',NULL),(62,41,23,'2024-07-06 08:52:06','Pending',NULL);
/*!40000 ALTER TABLE `BookRequests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Books`
--

DROP TABLE IF EXISTS `Books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Books` (
  `BookID` int NOT NULL AUTO_INCREMENT,
  `Title` varchar(255) NOT NULL,
  `Author` varchar(255) DEFAULT NULL,
  `Genre` varchar(255) DEFAULT NULL,
  `Quantity` int DEFAULT NULL,
  PRIMARY KEY (`BookID`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Books`
--

LOCK TABLES `Books` WRITE;
/*!40000 ALTER TABLE `Books` DISABLE KEYS */;
INSERT INTO `Books` VALUES (17,'hi','hello','waddup',6),(22,'increase','books','haha',8),(23,'decrease','books','haha',6),(24,'new','new','new',5);
/*!40000 ALTER TABLE `Books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `User`
--

DROP TABLE IF EXISTS `User`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `User` (
  `UserID` int NOT NULL AUTO_INCREMENT,
  `Username` varchar(255) NOT NULL,
  `Pass` varchar(255) NOT NULL,
  `Created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Role` varchar(255) NOT NULL DEFAULT 'user',
  `AdminRequest` enum('NoRequest','Pending','Accepted','Denied') DEFAULT 'NoRequest',
  PRIMARY KEY (`UserID`),
  UNIQUE KEY `unique_username` (`Username`)
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User`
--

LOCK TABLES `User` WRITE;
/*!40000 ALTER TABLE `User` DISABLE KEYS */;
INSERT INTO `User` VALUES (17,'admin','$2a$05$SutXVDSdiLxds./9oktZseaWfYdNn0PkWpSw4hVQ/gzt7VbPVfZ/W','2024-05-30 18:38:38','Admin','NoRequest'),(18,'helo','$2a$05$QUmKOy.vXUoFVAKopYYoJuQcFhoG.SWV1Ymvs6oCyNFNVDLNurVW6','2024-05-30 18:41:12','Admin','NoRequest'),(19,'beep','$2a$05$QGq/9LA.tc2RxJ/Shk/DWO/vvqQrukrEYQfqyZlahXkhKcBa/ZU0.','2024-06-04 20:04:45','Client','NoRequest'),(20,'ayush ','$2a$05$p.9CzmVO4VD5kHFQN7RIReRJ5gIgumKxJ7BAxEiQb2cmK7neVnkoC','2024-06-04 20:05:40','Client','NoRequest'),(21,'tejaswa','$2a$05$Zq4Zpeyaizs8mMJfpGOH8.MLBdHrs0JOZ64W/FIKR4bp5UdlgvqZm','2024-06-04 23:17:17','Admin','NoRequest'),(22,'black','$2a$05$XaQy1ZpLlkvEfceI56NyN.0NM4EGGVKPSrrDyJ9fo1ZQ.lFGRVbl2','2024-06-05 02:09:21','Client','NoRequest'),(23,'devesh','$2a$05$f2uGAeN8o5iofNNEiU2zde4Fz8D9gEEdT2N4sZ4RkX336pNerTEA.','2024-06-05 09:35:02','Admin','NoRequest'),(24,'dipak','$2a$05$0hNE5tGv1zmG7F31.dt6wuqCJFm1olvT2zmiX2jfEC.AXh0F7kDh6','2024-06-05 09:35:11','Client','NoRequest'),(25,'blah','$2a$05$gsueJVw6JVuWwE2AVY2hkuQU2GJiauaNQFZpBl0kgyal8Q3p4XP2C','2024-06-05 13:13:23','Admin','Accepted'),(26,'redragon','$2a$05$NfbjuuGwS87pBpzqQ0jb5./sRUB4mBvlhb/R4/LVyT4ILMGwg9lHe','2024-06-05 19:34:43','Client','Denied'),(27,'garv','$2a$05$5QrahNwzxNK9XzsOxSuo5OuYDcly8L1I0tXYD0VzpzdxfIQde7Lym','2024-06-06 11:13:04','Admin','Accepted'),(28,'gugu','$2a$05$2dLPBJcR2gHvyT71kKAYtesBjEhz.LfRbQ6adYsM0K8L4RxDFBfmO','2024-06-06 20:18:30','Client','NoRequest'),(29,'riki','$2a$05$mj0jmYSRNjx6ipfRc7FZb.COE7Zsm5TCVfa3cgeWHkLg56TSH9Opa','2024-06-08 08:59:46','Client','NoRequest'),(30,'ayush','$2a$05$eDfWlLshP61ebJX5T0Y2XOVIs3Wxq8gPXAnDxz7pVpAYmJkmsmvmu','2024-06-08 09:00:18','Client','NoRequest'),(31,'ishan','$2a$05$Uh3kC1uj33MRazfSxDZixONzae7/o3.sHPmSE9NxR0RFFDT3YVIJy','2024-06-08 12:39:36','Admin','Accepted'),(32,'sukhman','$2a$05$mrN7mFJQx8A7hjnkm/6hOuaFyAQzQtLYKOCNAaeYumM.jEadqJuI6','2024-06-08 23:07:53','Client','Pending'),(33,'blaze','$2a$05$jxxthvzD..sQcAFwitITReCIk1GxpOwAvp6lAmM3ikh4KWEtpX2tm','2024-06-08 23:08:20','Client','Pending'),(34,'','$2a$10$eSeqU213BOifgJz9oaSpkexeYNAi/LjbR/.qfG7mC2Kw//R8AM/uC','2024-06-29 20:51:25','Client','NoRequest'),(35,'new','$2a$10$v0SMEzvHEwjJ1tiuXtN9ieoFi4OMuhd2/KJAVWUBTgSnmVPJJm.mi','2024-06-29 20:53:14','Client','NoRequest'),(36,'boogi','$2a$10$9wjhWH81cWHAICmXI4GpV.v9IvElAIPOGaRmnDwBiBiGMM8WrcBWK','2024-07-01 17:49:00','Client','NoRequest'),(37,'blegh','$2a$10$uNxjkl7OTzMG0RF56TEzmeWHWtQEM/uDNzhf5nef5psVeOx1tQJOq','2024-07-01 18:18:03','Client','NoRequest'),(38,'mew','$2a$10$.UUQmeBb8mo.hg1S8347u.fHy8neotaiY.agjYY2rC3UATAjbbALi','2024-07-01 18:19:47','Client','NoRequest'),(41,'moww','$2a$10$ehm3W58ZEHkXpVztpsTxl.L7vIWtygcXBBRj4H9/eqp0w.6z3uzpq','2024-07-01 19:20:27','Client','NoRequest');
/*!40000 ALTER TABLE `User` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-07-07  4:22:11
