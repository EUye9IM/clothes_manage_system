-- MySQL dump 10.13  Distrib 5.5.62, for Win64 (AMD64)
--
-- Host: 192.168.80.200    Database: manage_system
-- ------------------------------------------------------
-- Server version	5.5.5-10.3.28-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `batch`
--

DROP TABLE IF EXISTS `batch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `batch` (
  `bt_id` int(11) NOT NULL AUTO_INCREMENT,
  `bt_u_id` int(11) NOT NULL,
  `bt_time` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`bt_id`),
  KEY `batch` (`bt_u_id`),
  CONSTRAINT `batch_ibfk_1` FOREIGN KEY (`bt_u_id`) REFERENCES `user` (`u_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `item`
--

DROP TABLE IF EXISTS `item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item` (
  `it_id` char(24) NOT NULL,
  `it_pd_id` int(11) NOT NULL,
  `it_bt_id` int(11) NOT NULL,
  PRIMARY KEY (`it_id`),
  KEY `it_pd_id` (`it_pd_id`),
  KEY `it_bt_id` (`it_bt_id`),
  CONSTRAINT `item_ibfk_1` FOREIGN KEY (`it_pd_id`) REFERENCES `product` (`pd_id`),
  CONSTRAINT `item_ibfk_2` FOREIGN KEY (`it_bt_id`) REFERENCES `batch` (`bt_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pattern`
--

DROP TABLE IF EXISTS `pattern`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pattern` (
  `pt_id` int(11) NOT NULL AUTO_INCREMENT,
  `pt_name` char(32) NOT NULL,
  `pt_brand` char(32) DEFAULT NULL,
  `pt_price` decimal(6,2) NOT NULL DEFAULT 0.00,
  PRIMARY KEY (`pt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product`
--

DROP TABLE IF EXISTS `product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product` (
  `pd_id` int(11) NOT NULL AUTO_INCREMENT,
  `pd_pt_id` int(11) NOT NULL,
  `pd_SKU` char(14) DEFAULT NULL,
  `pd_color` char(32) DEFAULT NULL,
  `pd_size` char(32) DEFAULT NULL,
  PRIMARY KEY (`pd_id`),
  KEY `pd_pt_id` (`pd_pt_id`),
  CONSTRAINT `product_ibfk_1` FOREIGN KEY (`pd_pt_id`) REFERENCES `pattern` (`pt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `u_id` int(11) NOT NULL AUTO_INCREMENT,
  `u_name` char(32) NOT NULL,
  `u_salt` char(8) NOT NULL,
  `u_pw` char(40) NOT NULL,
  `u_grant` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`u_id`),
  UNIQUE KEY `u_name` (`u_name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'manage_system'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-06-19  0:21:51