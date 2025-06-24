-- MySQL dump 10.13  Distrib 8.0.16, for Linux (x86_64)
--
-- Host: localhost    Database: aq3cms
-- ------------------------------------------------------
-- Server version	8.0.16

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8mb4 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `aq3cms_addonarticle`
--

DROP TABLE IF EXISTS `aq3cms_addonarticle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_addonarticle` (
  `aid` int(11) NOT NULL DEFAULT '0',
  `typeid` int(11) NOT NULL DEFAULT '0',
  `body` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `redirecturl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `templet` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `userip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`aid`),
  KEY `typeid` (`typeid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_addonarticle`
--

LOCK TABLES `aq3cms_addonarticle` WRITE;
/*!40000 ALTER TABLE `aq3cms_addonarticle` DISABLE KEYS */;
INSERT INTO `aq3cms_addonarticle` VALUES (1,0,'sdfsdfdf','','',''),(2,0,'sdfdsfdsfsdfds','','',''),(3,0,'dfgdgfdgfdg','','',''),(4,0,'dfgdfgfdg','','',''),(5,0,'fdgdfgdfgdfg','','',''),(6,0,'jjjjjjjjjjjjjjj','','',''),(7,0,'dsfdsfsdf','','',''),(10,0,'vbnbvnn','','',''),(11,0,'fgfdgfdg','','',''),(13,0,'fgfdg','','','');
/*!40000 ALTER TABLE `aq3cms_addonarticle` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_admin`
--

DROP TABLE IF EXISTS `aq3cms_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `rank` int(11) NOT NULL DEFAULT '10',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `last_login` datetime DEFAULT NULL,
  `last_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `login_count` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_admin`
--

LOCK TABLES `aq3cms_admin` WRITE;
/*!40000 ALTER TABLE `aq3cms_admin` DISABLE KEYS */;
INSERT INTO `aq3cms_admin` VALUES (1,'admin','0192023a7bbd73250516f069df18b500',10,'admin@example.com',1,NULL,'',0,'2025-05-30 10:07:39','2025-05-30 10:07:39');
/*!40000 ALTER TABLE `aq3cms_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_archives`
--

DROP TABLE IF EXISTS `aq3cms_archives`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_archives` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `typeid` int(11) NOT NULL DEFAULT '0',
  `typeid2` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0',
  `sortrank` int(11) NOT NULL DEFAULT '0',
  `flag` set('c','h','p','f','s','j','a','b') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ismake` tinyint(1) NOT NULL DEFAULT '0',
  `channel` int(11) NOT NULL DEFAULT '1',
  `arcrank` int(11) NOT NULL DEFAULT '0',
  `click` int(11) NOT NULL DEFAULT '0',
  `money` int(11) NOT NULL DEFAULT '0',
  `title` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `shorttitle` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `color` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `writer` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `source` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `litpic` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `pubdate` int(11) NOT NULL DEFAULT '0',
  `senddate` int(11) NOT NULL DEFAULT '0',
  `mid` int(11) NOT NULL DEFAULT '0',
  `keywords` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `lastpost` int(11) NOT NULL DEFAULT '0',
  `scores` int(11) NOT NULL DEFAULT '0',
  `goodpost` int(11) NOT NULL DEFAULT '0',
  `badpost` int(11) NOT NULL DEFAULT '0',
  `voteid` int(11) NOT NULL DEFAULT '0',
  `notpost` tinyint(1) NOT NULL DEFAULT '0',
  `description` varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `filename` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `dutyadmin` int(11) NOT NULL DEFAULT '0',
  `tackid` int(11) NOT NULL DEFAULT '0',
  `mtype` int(11) NOT NULL DEFAULT '0',
  `weight` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `typeid` (`typeid`),
  KEY `sortrank` (`sortrank`),
  KEY `pubdate` (`pubdate`),
  KEY `click` (`click`),
  KEY `arcrank` (`arcrank`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_archives`
--

LOCK TABLES `aq3cms_archives` WRITE;
/*!40000 ALTER TABLE `aq3cms_archives` DISABLE KEYS */;
INSERT INTO `aq3cms_archives` VALUES (1,2,'0',0,'',0,1,0,2,0,'aaaaaaaaaaa','','#000000','admin','','',1748602677,1748602677,0,'',0,0,0,0,0,0,'sdfsdf','',0,0,0,0),(2,1,'0',0,'',0,1,0,1,0,'fsdfsdfdsf','','#000000','admin','','',1748866888,1748866888,0,'',0,0,0,0,0,0,'sdfdsf','',0,0,0,0),(3,3,'0',0,'',0,1,0,1,0,'dfgdfgdfg','','#000000','admin','','',1748866909,1748866909,0,'',0,0,0,0,0,0,'','',0,0,0,0),(4,4,'0',0,'',0,1,0,0,0,'refdgfdg','','#000000','admin','','',1748867035,1748867035,0,'',0,0,0,0,0,0,'','',0,0,0,0),(5,4,'0',0,'',0,1,0,0,0,'dfgdfgdfg555','','#000000','admin','','',1748867066,1748867066,0,'',0,0,0,0,0,0,'','',0,0,0,0),(6,4,'0',0,'',0,1,0,1,0,'jjghjhgjhg','','#000000','admin','','',1748868878,1748868878,0,'',0,0,0,0,0,0,'','',0,0,0,0),(7,2,'0',0,'',0,1,0,0,0,'dsfdfsdf888','','#000000','admin','','',1748869643,1748869643,0,'',0,0,0,0,0,0,'','',0,0,0,0),(8,5,'0',0,NULL,0,1,0,0,0,'1','','','admin','','',1748898000,1748898000,0,',',0,0,0,0,0,0,'','',0,0,0,0),(9,5,'0',0,NULL,0,1,0,1,0,'2','','','admin','','',1748898300,1748898300,0,',',0,0,0,0,0,0,'','',0,0,0,0),(10,5,'0',0,'',0,1,0,1,0,'cvbbnnbvn','','#000000','admin','','',1748872346,1748872346,0,'',0,0,0,0,0,0,'','',0,0,0,0),(11,5,'0',0,'',0,1,0,0,0,'fddgdff111','','#000000','admin','','',1748872371,1748872371,0,'',0,0,0,0,0,0,'','',0,0,0,0),(12,6,'0',0,NULL,0,1,0,0,0,'1','','','admin','','',1748901000,1748901000,0,',,',0,0,0,0,0,0,'','',0,0,0,0),(13,5,'0',0,'',0,1,0,1,0,'dfgdfg666','','#000000','admin','','',1748872399,1748872399,0,'',0,0,0,0,0,0,'','',0,0,0,0),(14,6,'0',0,NULL,0,1,0,0,0,'2','','','admin','','',1748901300,1748901300,0,',,',0,0,0,0,0,0,'','',0,0,0,0),(15,7,'0',0,NULL,0,1,0,0,0,'','','','admin','','',1748904000,1748904000,0,',AI,',0,0,0,0,0,0,'','',0,0,0,0);
/*!40000 ALTER TABLE `aq3cms_archives` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_arctype`
--

DROP TABLE IF EXISTS `aq3cms_arctype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_arctype` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `reid` int(11) NOT NULL DEFAULT '0',
  `topid` int(11) NOT NULL DEFAULT '0',
  `sortrank` int(11) NOT NULL DEFAULT '50',
  `typename` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `typedir` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `isdefault` tinyint(1) NOT NULL DEFAULT '0',
  `defaultname` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'index.html',
  `issend` tinyint(1) NOT NULL DEFAULT '1',
  `channeltype` int(11) NOT NULL DEFAULT '1',
  `maxpage` int(11) NOT NULL DEFAULT '-1',
  `ispart` tinyint(1) NOT NULL DEFAULT '0',
  `corank` int(11) NOT NULL DEFAULT '0',
  `tempindex` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `templist` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `temparticle` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `namerule` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `namerule2` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `modname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `keywords` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `seotitle` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `moresite` tinyint(1) NOT NULL DEFAULT '0',
  `sitepath` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `siteurl` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `ishidden` tinyint(1) NOT NULL DEFAULT '0',
  `cross` tinyint(1) NOT NULL DEFAULT '0',
  `crossid` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `smalltypes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `reid` (`reid`),
  KEY `topid` (`topid`),
  KEY `sortrank` (`sortrank`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_arctype`
--

LOCK TABLES `aq3cms_arctype` WRITE;
/*!40000 ALTER TABLE `aq3cms_arctype` DISABLE KEYS */;
INSERT INTO `aq3cms_arctype` VALUES (1,0,0,50,'栏目名323','aaa',1,'index.html',1,1,-1,0,0,'index.htm','','','','','','ç½‘ç«™é¦–é¡µ','','',0,'','',0,0,'0',NULL,NULL),(2,0,0,50,'栏目名bbb','news',0,'index.html',1,1,-1,0,0,'','','','/news/{Y}/{M}{D}/{aid}.html','/news/list_{tid}_{page}.html','','æ–°é—»ä¸­å¿ƒ','æ–°é—»','æ–°é—»ä¸­å¿ƒ',0,'','',0,0,'0',NULL,NULL),(3,0,0,50,'栏目名ccc','product',0,'index.html',1,1,-1,0,0,'','','','/product/{Y}/{M}{D}/{aid}.html','/product/list_{tid}_{page}.html','','äº§å“å±•ç¤º','äº§å“','äº§å“å±•ç¤º',0,'','',0,0,'0',NULL,NULL),(4,0,0,50,'栏目名sdf','about',0,'index.html',1,1,-1,1,0,'','','','','','','å…³äºŽæˆ‘ä»¬','å…³äºŽ','å…³äºŽæˆ‘ä»¬',0,'','',0,0,'0',NULL,NULL),(5,0,0,50,'products','products',0,'index.html',1,1,-1,0,0,'','','','','','','products','products','',0,'','',0,0,'0',NULL,NULL),(6,5,0,60,'','featured',0,'index.html',1,1,-1,0,0,'','','','','','','',',','',0,'','',0,0,NULL,NULL,NULL),(7,1,0,70,'','tech',0,'index.html',1,1,-1,0,0,'','','','','','','',',','',0,'','',0,0,NULL,NULL,NULL);
/*!40000 ALTER TABLE `aq3cms_arctype` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_feedback`
--

DROP TABLE IF EXISTS `aq3cms_feedback`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_feedback` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'è¯„è®ºID',
  `aid` int(11) NOT NULL DEFAULT '0' COMMENT 'æ–‡ç« ID',
  `typeid` int(11) NOT NULL DEFAULT '0' COMMENT 'æ ç›®ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'ç”¨æˆ·å',
  `mid` int(11) NOT NULL DEFAULT '0' COMMENT 'ä¼šå‘˜ID',
  `ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'IPåœ°å€',
  `ischeck` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'å®¡æ ¸çŠ¶æ€ï¼š0=å¾…å®¡æ ¸ï¼Œ1=å·²å®¡æ ¸ï¼Œ-1=å·²æ‹’ç»',
  `dtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'è¯„è®ºæ—¶é—´',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è¯„è®ºå†…å®¹',
  `parentid` int(11) NOT NULL DEFAULT '0' COMMENT 'çˆ¶è¯„è®ºIDï¼ˆå›žå¤åŠŸèƒ½ï¼‰',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT 'è¯„åˆ†',
  `goodcount` int(11) NOT NULL DEFAULT '0' COMMENT 'ç‚¹èµžæ•°',
  `badcount` int(11) NOT NULL DEFAULT '0' COMMENT 'è¸©æ•°',
  `userface` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'ç”¨æˆ·å¤´åƒ',
  `channeltype` tinyint(1) NOT NULL DEFAULT '1' COMMENT 'é¢‘é“ç±»åž‹ï¼š1=æ–‡ç« ',
  PRIMARY KEY (`id`),
  KEY `idx_aid` (`aid`),
  KEY `idx_typeid` (`typeid`),
  KEY `idx_mid` (`mid`),
  KEY `idx_ischeck` (`ischeck`),
  KEY `idx_dtime` (`dtime`),
  KEY `idx_parentid` (`parentid`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è¯„è®ºåé¦ˆè¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_feedback`
--

LOCK TABLES `aq3cms_feedback` WRITE;
/*!40000 ALTER TABLE `aq3cms_feedback` DISABLE KEYS */;
INSERT INTO `aq3cms_feedback` VALUES (1,1,1,'å¼ ä¸‰',1,'192.168.1.100',1,'2024-01-15 10:30:00','è¿™ç¯‡æ–‡ç« å†™å¾—å¾ˆå¥½ï¼Œå†…å®¹è¯¦å®žï¼Œå¯¹æˆ‘å¾ˆæœ‰å¸®åŠ©ï¼',0,5,12,1,'/uploads/face/default.jpg',1),(2,1,1,'æŽå››',2,'192.168.1.101',1,'2024-01-15 11:15:00','åŒæ„æ¥¼ä¸Šçš„è§‚ç‚¹ï¼Œä½œè€…çš„åˆ†æžå¾ˆåˆ°ä½ã€‚',1,4,8,0,'/uploads/face/default.jpg',1),(3,1,1,'çŽ‹äº”',0,'192.168.1.102',1,'2024-01-15 14:20:00','æœ‰äº›åœ°æ–¹è¿˜å¯ä»¥æ›´æ·±å…¥ä¸€äº›ï¼ŒæœŸå¾…ä½œè€…çš„åŽç»­æ–‡ç« ã€‚',0,4,3,0,'/uploads/face/default.jpg',1),(4,2,1,'èµµå…­',3,'192.168.1.103',1,'2024-01-16 09:45:00','æŠ€æœ¯æ–‡ç« å°±åº”è¯¥è¿™æ ·å†™ï¼Œæ¸…æ™°æ˜Žäº†ï¼Œå®žç”¨æ€§å¼ºã€‚',0,5,15,0,'/uploads/face/default.jpg',1),(5,2,1,'é’±ä¸ƒ',0,'192.168.1.104',-1,'2024-01-16 16:30:00','ä»£ç ç¤ºä¾‹å¾ˆæ£’ï¼Œå·²ç»åœ¨é¡¹ç›®ä¸­åº”ç”¨äº†ï¼Œæ•ˆæžœä¸é”™ï¼',0,5,6,0,'/uploads/face/default.jpg',1),(6,3,2,'å­™å…«',4,'192.168.1.105',1,'2024-01-17 08:20:00','è¿™ä¸ªè§£å†³æ–¹æ¡ˆç¡®å®žæœ‰æ•ˆï¼Œæ„Ÿè°¢åˆ†äº«ï¼',0,4,9,1,'/uploads/face/default.jpg',1),(7,3,2,'å‘¨ä¹',0,'192.168.1.106',-1,'2024-01-17 12:10:00','å†…å®¹éœ€è¦æ”¹è¿›ï¼Œå¸Œæœ›èƒ½æä¾›æ›´å¤šå®žç”¨çš„ä¾‹å­ã€‚',0,2,1,3,'/uploads/face/default.jpg',1),(8,1,1,'å´å',5,'192.168.1.107',1,'2024-01-17 15:45:00','@æŽå›› æˆ‘ä¹Ÿæœ‰åŒæ ·çš„æ„Ÿå—ï¼Œä½œè€…çš„æ€è·¯å¾ˆæ¸…æ™°ã€‚',2,4,5,0,'/uploads/face/default.jpg',1),(10,4,3,'çŽ‹åäºŒ',6,'192.168.1.109',1,'2024-01-18 14:30:00','æ–‡æ¡£å†™å¾—å¾ˆè¯¦ç»†ï¼ŒæŒ‰ç…§æ­¥éª¤æ“ä½œæˆåŠŸäº†ï¼',0,5,11,0,'/uploads/face/default.jpg',1);
/*!40000 ALTER TABLE `aq3cms_feedback` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_flink`
--

DROP TABLE IF EXISTS `aq3cms_flink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_flink` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sortrank` smallint(6) NOT NULL DEFAULT '0' COMMENT 'æŽ’åº',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT 'é“¾æŽ¥åœ°å€',
  `webname` varchar(60) NOT NULL DEFAULT '' COMMENT 'ç½‘ç«™åç§°',
  `msg` text COMMENT 'ç½‘ç«™æè¿°',
  `email` varchar(60) NOT NULL DEFAULT '' COMMENT 'è”ç³»é‚®ç®±',
  `typeid` smallint(6) NOT NULL DEFAULT '0' COMMENT 'åˆ†ç±»ID',
  `logo` varchar(255) NOT NULL DEFAULT '' COMMENT 'LOGOåœ°å€',
  `ischeck` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'æ˜¯å¦å®¡æ ¸é€šè¿‡',
  `dtime` int(11) NOT NULL DEFAULT '0' COMMENT 'æ·»åŠ æ—¶é—´',
  PRIMARY KEY (`id`),
  KEY `sortrank` (`sortrank`),
  KEY `typeid` (`typeid`),
  KEY `ischeck` (`ischeck`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='å‹æƒ…é“¾æŽ¥è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_flink`
--

LOCK TABLES `aq3cms_flink` WRITE;
/*!40000 ALTER TABLE `aq3cms_flink` DISABLE KEYS */;
INSERT INTO `aq3cms_flink` VALUES (1,1,'https://golang.org','Go å®˜æ–¹ç½‘ç«™','Go ç¼–ç¨‹è¯­è¨€å®˜æ–¹ç½‘ç«™','admin@golang.org',1,'',1,1748676863),(2,2,'https://github.com','GitHub','å…¨çƒæœ€å¤§çš„ä»£ç æ‰˜ç®¡å¹³å°','support@github.com',1,'',1,1748676863),(3,3,'https://getbootstrap.com','Bootstrap','æµè¡Œçš„å‰ç«¯æ¡†æž¶','support@getbootstrap.com',1,'',1,1748676863),(4,4,'https://jquery.com','jQuery','JavaScript åº“','support@jquery.com',1,'',1,1748676863),(5,5,'https://mysql.com','MySQL','æµè¡Œçš„å…³ç³»åž‹æ•°æ®åº“','support@mysql.com',1,'',1,1748676863);
/*!40000 ALTER TABLE `aq3cms_flink` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_member`
--

DROP TABLE IF EXISTS `aq3cms_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_member` (
  `mid` int(11) NOT NULL AUTO_INCREMENT,
  `mtype` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `userid` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `uname` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `sex` enum('ç”·','å¥³','ä¿å¯†') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ä¿å¯†',
  `rank` int(11) NOT NULL DEFAULT '10',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `scores` int(11) NOT NULL DEFAULT '0',
  `matt` tinyint(1) NOT NULL DEFAULT '0',
  `spacesta` tinyint(1) NOT NULL DEFAULT '0',
  `face` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `safequestion` tinyint(1) NOT NULL DEFAULT '0',
  `safeanswer` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `jointime` int(11) NOT NULL DEFAULT '0',
  `joinip` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `logintime` int(11) NOT NULL DEFAULT '0',
  `loginip` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`mid`),
  KEY `userid` (`userid`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_member`
--

LOCK TABLES `aq3cms_member` WRITE;
/*!40000 ALTER TABLE `aq3cms_member` DISABLE KEYS */;
INSERT INTO `aq3cms_member` VALUES (1,'1','testuser1','e10adc3949ba59abbe56e057f20f883e','','ç”·',10,0.00,'test1@example.com',100,0,0,'',0,'',1705312800,'192.168.1.100',1705764600,'192.168.1.100'),(2,'2','testuser2','e10adc3949ba59abbe56e057f20f883e','','',10,50.00,'test2@example.com',200,0,0,'',0,'',1705402800,'192.168.1.101',1705854600,'192.168.1.101'),(3,'1','testuser3','e10adc3949ba59abbe56e057f20f883e','','ç”·',10,25.00,'test3@example.com',150,0,0,'',0,'',1705492800,'192.168.1.102',1705944600,'192.168.1.102'),(4,'3','testuser4','e10adc3949ba59abbe56e057f20f883e','','ç”·',10,100.00,'test4@example.com',300,0,0,'',0,'',1705582800,'192.168.1.103',1706034600,'192.168.1.103'),(5,'2','testuser5','e10adc3949ba59abbe56e057f20f883e','','ç”·',10,75.00,'test5@example.com',250,0,0,'',0,'',1705672800,'192.168.1.104',1706124600,'192.168.1.104'),(6,'1','newuser1','e10adc3949ba59abbe56e057f20f883e','','ç”·',10,0.00,'new1@example.com',50,0,0,'',0,'',1748603674,'192.168.1.105',1748603674,'192.168.1.105'),(7,'1','newuser2','e10adc3949ba59abbe56e057f20f883e','','ç”·',10,0.00,'new2@example.com',60,0,0,'',0,'',1748603674,'192.168.1.106',1748603674,'192.168.1.106'),(8,'1','dfsdfwerwer','0567212991e29fa44ca3a6189a911a5f','','å¥³',10,0.00,'dfdfsf@df.com',0,0,0,'',0,'',1748653232,'127.0.0.1:55596',1748653232,''),(9,'1','adminsdfadfddf','0567212991e29fa44ca3a6189a911a5f','','å¥³',10,0.00,'dfdfdfsf@df.com',0,0,0,'',0,'',1748653249,'127.0.0.1:55596',1748653249,'');
/*!40000 ALTER TABLE `aq3cms_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_member_type`
--

DROP TABLE IF EXISTS `aq3cms_member_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_member_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `typename` varchar(50) NOT NULL DEFAULT '',
  `description` text,
  `rank` int(11) NOT NULL DEFAULT '10',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00',
  `scores` int(11) NOT NULL DEFAULT '0',
  `purviews` text,
  PRIMARY KEY (`id`),
  KEY `idx_rank` (`rank`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ä¼šå‘˜ç±»åž‹è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_member_type`
--

LOCK TABLES `aq3cms_member_type` WRITE;
/*!40000 ALTER TABLE `aq3cms_member_type` DISABLE KEYS */;
INSERT INTO `aq3cms_member_type` VALUES (1,'aaa','aaa1',10,0.00,0,NULL),(2,'VIP','vip1',20,100.00,1000,NULL),(3,'bbb','bbb1',30,500.00,5000,NULL),(4,'ccc','ccc1',40,1000.00,10000,NULL);
/*!40000 ALTER TABLE `aq3cms_member_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_sysconfig`
--

DROP TABLE IF EXISTS `aq3cms_sysconfig`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_sysconfig` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `varname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `info` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `groupid` tinyint(1) NOT NULL DEFAULT '1',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'string',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `varname` (`varname`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_sysconfig`
--

LOCK TABLES `aq3cms_sysconfig` WRITE;
/*!40000 ALTER TABLE `aq3cms_sysconfig` DISABLE KEYS */;
INSERT INTO `aq3cms_sysconfig` VALUES (1,'cfg_webname','ç½‘ç«™åç§°',1,'string','aq3cmså†…å®¹ç®¡ç†ç³»ç»Ÿ'),(2,'cfg_weburl','ç½‘ç«™åœ°å€',1,'string','http://localhost'),(3,'cfg_keywords','ç½‘ç«™å…³é”®è¯',1,'string','aq3cms,å†…å®¹ç®¡ç†ç³»ç»Ÿ,Goè¯­è¨€'),(4,'cfg_description','ç½‘ç«™æè¿°',1,'string','aq3cmsæ˜¯åŸºäºŽGoè¯­è¨€å¼€å‘çš„çŽ°ä»£åŒ–å†…å®¹ç®¡ç†ç³»ç»Ÿ'),(5,'cfg_icp','ICPå¤‡æ¡ˆå·',1,'string',''),(6,'cfg_statcode','ç»Ÿè®¡ä»£ç ',1,'text',''),(7,'cfg_copyright','ç‰ˆæƒä¿¡æ¯',1,'string','Copyright Â© 2024 aq3cms. All rights reserved.');
/*!40000 ALTER TABLE `aq3cms_sysconfig` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_tagindex`
--

DROP TABLE IF EXISTS `aq3cms_tagindex`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_tagindex` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tag` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `count` int(11) NOT NULL DEFAULT '0',
  `rank` int(11) NOT NULL DEFAULT '0',
  `ishot` tinyint(1) NOT NULL DEFAULT '0',
  `addtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastuse` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `tagpinyin` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `tag` (`tag`),
  KEY `rank` (`rank`),
  KEY `ishot` (`ishot`),
  KEY `count` (`count`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_tagindex`
--

LOCK TABLES `aq3cms_tagindex` WRITE;
/*!40000 ALTER TABLE `aq3cms_tagindex` DISABLE KEYS */;
INSERT INTO `aq3cms_tagindex` VALUES (1,'建议使',5,10,1,'2025-05-30 10:08:06','2025-05-30 10:08:06','jishu'),(2,'fdgfdg',8,20,1,'2025-05-30 10:08:06','2025-05-30 10:08:06','xinwen'),(3,'建议使dfdf2',3,15,0,'2025-05-30 10:08:06','2025-05-30 10:08:06','chanpin'),(4,'建议使sdsd',6,12,1,'2025-05-30 10:08:06','2025-05-30 10:08:06','jiaocheng'),(6,'sdfsdfdsf',0,0,0,'2025-05-30 18:08:17','2025-05-30 18:08:17','sdfsdfdsf'),(11,'å®‰å…¨',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(12,'æ¨¡æ¿',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(13,'æ’ä»¶',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(14,'æ•°æ®åº“',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(15,'æ€§èƒ½',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(16,'ç¼“å­˜',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(17,'CMS',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(18,'aq3cms',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(19,'Golang',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(20,'Webå¼€å‘',1,0,0,'2025-06-02 18:24:40','2025-06-02 18:24:40',''),(21,'dfgdg',1,0,0,'2025-06-02 20:24:26','2025-06-02 20:24:26','');
/*!40000 ALTER TABLE `aq3cms_tagindex` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aq3cms_taglist`
--

DROP TABLE IF EXISTS `aq3cms_taglist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `aq3cms_taglist` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) NOT NULL DEFAULT '0',
  `tag` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `addtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `aid_tag` (`aid`,`tag`),
  KEY `aid` (`aid`),
  KEY `tag` (`tag`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aq3cms_taglist`
--

LOCK TABLES `aq3cms_taglist` WRITE;
/*!40000 ALTER TABLE `aq3cms_taglist` DISABLE KEYS */;
INSERT INTO `aq3cms_taglist` VALUES (11,1,'å®‰å…¨','2025-06-02 10:24:40'),(12,1,'æ¨¡æ¿','2025-06-02 10:24:40'),(13,1,'æ’ä»¶','2025-06-02 10:24:40'),(14,1,'æ•°æ®åº“','2025-06-02 10:24:40'),(15,1,'æ€§èƒ½','2025-06-02 10:24:40'),(16,1,'ç¼“å­˜','2025-06-02 10:24:40'),(17,1,'CMS','2025-06-02 10:24:40'),(18,1,'aq3cms','2025-06-02 10:24:40'),(19,1,'Golang','2025-06-02 10:24:40'),(20,1,'Webå¼€å‘','2025-06-02 10:24:40'),(21,5,'dfgdg','2025-06-02 12:24:26');
/*!40000 ALTER TABLE `aq3cms_taglist` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-13  7:30:59
