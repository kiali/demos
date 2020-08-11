# Initialize a mysql db with a 'test' db and be able test productpage with it.
# mysql -h 127.0.0.1 -ppassword < mysqldb-init.sql

CREATE DATABASE test;
USE test;

CREATE TABLE `cities`
(
    `cityId` INT NOT NULL,
    `city`   TEXT,
    `lat`    TEXT,
    `lng`    TEXT,
    PRIMARY KEY (`cityId`)
);

CREATE TABLE `cars`
(
    `carId`    INT NOT NULL,
    `cityId`   INT,
    `carModel` TEXT,
    `price`    DECIMAL,
    PRIMARY KEY (`carId`),
    FOREIGN KEY (`cityId`) REFERENCES cities (cityId)
);

CREATE TABLE `flights`
(
    `flightId` INT NOT NULL,
    `cityId`   INT,
    `airline`  TEXT,
    `price`    DECIMAL,
    PRIMARY KEY (`flightId`),
    FOREIGN KEY (`cityId`) REFERENCES cities (cityId)
);

CREATE TABLE `hotels`
(
    `hotelId` INT NOT NULL,
    `cityId`  INT,
    `hotel`   TEXT,
    `price`   DECIMAL,
    PRIMARY KEY (`hotelId`),
    FOREIGN KEY (`cityId`) REFERENCES cities (cityId)
);

CREATE TABLE `insurances`
(
    `insuranceId` INT NOT NULL,
    `cityId`      INT,
    `company`     TEXT,
    `price`       DECIMAL,
    PRIMARY KEY (`insuranceId`),
    FOREIGN KEY (`cityId`) REFERENCES cities (cityId)
);

INSERT INTO cities (cityId, city, lat, lng) VALUES(1, 'Amsterdam', '52.3500', '4.9166');
INSERT INTO cities (cityId, city, lat, lng) VALUES(2, 'Andorra', '42.5000', '1.5165');
INSERT INTO cities (cityId, city, lat, lng) VALUES(3, 'Athens', '37.9833', '23.7333');
INSERT INTO cities (cityId, city, lat, lng) VALUES(4, 'Belgrade', '44.8186', '20.4680');
INSERT INTO cities (cityId, city, lat, lng) VALUES(5, 'Berlin', '52.5218', '13.4015');
INSERT INTO cities (cityId, city, lat, lng) VALUES(6, 'Bern', '46.9167', '7.4670');
INSERT INTO cities (cityId, city, lat, lng) VALUES(7, 'Bratislava', '48.1500', '17.1170');
INSERT INTO cities (cityId, city, lat, lng) VALUES(8, 'Brussels', '50.8333', '4.3333');
INSERT INTO cities (cityId, city, lat, lng) VALUES(9, 'Bucharest', '44.4334', '26.0999');
INSERT INTO cities (cityId, city, lat, lng) VALUES(10, 'Budapest', '47.5000', '19.0833');
INSERT INTO cities (cityId, city, lat, lng) VALUES(11, 'Chisinau', '47.0050', '28.8577');
INSERT INTO cities (cityId, city, lat, lng) VALUES(12, 'Copenhagen', '55.6786', '12.5635');
INSERT INTO cities (cityId, city, lat, lng) VALUES(13, 'Dublin', '53.3331', '-6.2489');
INSERT INTO cities (cityId, city, lat, lng) VALUES(14, 'Helsinki', '60.1756', '24.9341');
INSERT INTO cities (cityId, city, lat, lng) VALUES(15, 'Kiev', '50.473782', '30.516237');
INSERT INTO cities (cityId, city, lat, lng) VALUES(16, 'Lisbon', '38.7227', '-9.1449');
INSERT INTO cities (cityId, city, lat, lng) VALUES(17, 'Ljubljana', '46.0553', '14.5150');
INSERT INTO cities (cityId, city, lat, lng) VALUES(18, 'London', '51.5000', '-0.1167');
INSERT INTO cities (cityId, city, lat, lng) VALUES(19, 'Luxembourg', '49.6117', '6.1300');
INSERT INTO cities (cityId, city, lat, lng) VALUES(20, 'Madrid', '40.4000', '-3.6834');
INSERT INTO cities (cityId, city, lat, lng) VALUES(21, 'Minsk', '53.9000', '27.5666');
INSERT INTO cities (cityId, city, lat, lng) VALUES(22, 'Monaco', '43.7396', '7.4069');
INSERT INTO cities (cityId, city, lat, lng) VALUES(23, 'Moscow', '55.7522', '37.6155');
INSERT INTO cities (cityId, city, lat, lng) VALUES(24, 'Nicosia', '35.1667', '33.3666');
INSERT INTO cities (cityId, city, lat, lng) VALUES(25, 'Nuuk', '64.1983', '-51.7327');
INSERT INTO cities (cityId, city, lat, lng) VALUES(26, 'Oslo', '59.9167', '10.7500');
INSERT INTO cities (cityId, city, lat, lng) VALUES(27, 'Paris', '48.8667', '2.3333');
INSERT INTO cities (cityId, city, lat, lng) VALUES(28, 'Podgorica', '42.4660', '19.2663');
INSERT INTO cities (cityId, city, lat, lng) VALUES(29, 'Prague', '50.0833', '14.4660');
INSERT INTO cities (cityId, city, lat, lng) VALUES(30, 'Reykjavik', '64.1500', '-21.9500');
INSERT INTO cities (cityId, city, lat, lng) VALUES(31, 'Riga', '56.9500', '24.1000');
INSERT INTO cities (cityId, city, lat, lng) VALUES(32, 'Rome', '41.8960', '12.4833');
INSERT INTO cities (cityId, city, lat, lng) VALUES(33, 'San Marino', '43.9172', '12.4667');
INSERT INTO cities (cityId, city, lat, lng) VALUES(34, 'Sarajevo', '43.8500', '18.3830');
INSERT INTO cities (cityId, city, lat, lng) VALUES(35, 'Skopje', '42.0000', '21.4335');
INSERT INTO cities (cityId, city, lat, lng) VALUES(36, 'Sofia', '42.6833', '23.3167');
INSERT INTO cities (cityId, city, lat, lng) VALUES(37, 'Stockholm', '59.3508', '18.0973');
INSERT INTO cities (cityId, city, lat, lng) VALUES(38, 'Tallinn', '59.4339', '24.7280');
INSERT INTO cities (cityId, city, lat, lng) VALUES(39, 'Tirana', '41.3275', '19.8189');
INSERT INTO cities (cityId, city, lat, lng) VALUES(40, 'Vaduz', '47.1337', '9.5167');
INSERT INTO cities (cityId, city, lat, lng) VALUES(41, 'Valletta', '35.8997', '14.5147');
INSERT INTO cities (cityId, city, lat, lng) VALUES(42, 'Vienna', '48.2000', '16.3666');
INSERT INTO cities (cityId, city, lat, lng) VALUES(43, 'Vilnius', '54.6834', '25.3166');
INSERT INTO cities (cityId, city, lat, lng) VALUES(44, 'Warsaw', '52.2500', '21.0000');
INSERT INTO cities (cityId, city, lat, lng) VALUES(45, 'Zagreb', '45.8000', '16.0000');

INSERT INTO cars (carId, cityId, carModel, price) VALUES(1, 1, 'Sports Car', 1005);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(2, 1, 'Economy Car', 302);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(3, 2, 'Sports Car', 1010);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(4, 2, 'Economy Car', 304);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(5, 3, 'Sports Car', 1015);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(6, 3, 'Economy Car', 306);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(7, 4, 'Sports Car', 1020);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(8, 4, 'Economy Car', 308);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(9, 5, 'Sports Car', 1025);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(10, 5, 'Economy Car', 310);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(11, 6, 'Sports Car', 1030);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(12, 6, 'Economy Car', 312);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(13, 7, 'Sports Car', 1035);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(14, 7, 'Economy Car', 314);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(15, 8, 'Sports Car', 1040);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(16, 8, 'Economy Car', 316);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(17, 9, 'Sports Car', 1045);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(18, 9, 'Economy Car', 318);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(19, 10, 'Sports Car', 1050);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(20, 10, 'Economy Car', 320);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(21, 11, 'Sports Car', 1055);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(22, 11, 'Economy Car', 322);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(23, 12, 'Sports Car', 1060);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(24, 12, 'Economy Car', 324);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(25, 13, 'Sports Car', 1065);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(26, 13, 'Economy Car', 326);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(27, 14, 'Sports Car', 1070);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(28, 14, 'Economy Car', 328);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(29, 15, 'Sports Car', 1075);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(30, 15, 'Economy Car', 330);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(31, 16, 'Sports Car', 1080);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(32, 16, 'Economy Car', 332);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(33, 17, 'Sports Car', 1085);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(34, 17, 'Economy Car', 334);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(35, 18, 'Sports Car', 1090);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(36, 18, 'Economy Car', 336);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(37, 19, 'Sports Car', 1095);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(38, 19, 'Economy Car', 338);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(39, 20, 'Sports Car', 1100);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(40, 20, 'Economy Car', 340);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(41, 21, 'Sports Car', 1105);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(42, 21, 'Economy Car', 342);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(43, 22, 'Sports Car', 1110);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(44, 22, 'Economy Car', 344);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(45, 23, 'Sports Car', 1115);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(46, 23, 'Economy Car', 346);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(47, 24, 'Sports Car', 1120);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(48, 24, 'Economy Car', 348);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(49, 25, 'Sports Car', 1125);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(50, 25, 'Economy Car', 350);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(51, 26, 'Sports Car', 1130);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(52, 26, 'Economy Car', 352);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(53, 27, 'Sports Car', 1135);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(54, 27, 'Economy Car', 354);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(55, 28, 'Sports Car', 1140);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(56, 28, 'Economy Car', 356);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(57, 29, 'Sports Car', 1145);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(58, 29, 'Economy Car', 358);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(59, 30, 'Sports Car', 1150);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(60, 30, 'Economy Car', 360);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(61, 31, 'Sports Car', 1155);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(62, 31, 'Economy Car', 362);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(63, 32, 'Sports Car', 1160);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(64, 32, 'Economy Car', 364);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(65, 33, 'Sports Car', 1165);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(66, 33, 'Economy Car', 366);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(67, 34, 'Sports Car', 1170);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(68, 34, 'Economy Car', 368);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(69, 35, 'Sports Car', 1175);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(70, 35, 'Economy Car', 370);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(71, 36, 'Sports Car', 1180);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(72, 36, 'Economy Car', 372);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(73, 37, 'Sports Car', 1185);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(74, 37, 'Economy Car', 374);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(75, 38, 'Sports Car', 1190);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(76, 38, 'Economy Car', 376);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(77, 39, 'Sports Car', 1195);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(78, 39, 'Economy Car', 378);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(79, 40, 'Sports Car', 1200);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(80, 40, 'Economy Car', 380);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(81, 41, 'Sports Car', 1205);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(82, 41, 'Economy Car', 382);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(83, 42, 'Sports Car', 1210);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(84, 42, 'Economy Car', 384);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(85, 43, 'Sports Car', 1215);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(86, 43, 'Economy Car', 386);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(87, 44, 'Sports Car', 1220);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(88, 44, 'Economy Car', 388);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(89, 45, 'Sports Car', 1225);
INSERT INTO cars (carId, cityId, carModel, price) VALUES(90, 45, 'Economy Car', 390);

INSERT INTO flights (flightId, cityId, airline, price) VALUES(1, 1, 'Red Airlines', 1001);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(2, 1, 'Blue Airlines', 351);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(3, 1, 'Green Airlines', 301);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(4, 2, 'Red Airlines', 1002);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(5, 2, 'Blue Airlines', 352);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(6, 2, 'Green Airlines', 302);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(7, 3, 'Red Airlines', 1003);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(8, 3, 'Blue Airlines', 353);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(9, 3, 'Green Airlines', 303);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(10, 4, 'Red Airlines', 1004);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(11, 4, 'Blue Airlines', 354);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(12, 4, 'Green Airlines', 304);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(13, 5, 'Red Airlines', 1005);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(14, 5, 'Blue Airlines', 355);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(15, 5, 'Green Airlines', 305);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(16, 6, 'Red Airlines', 1006);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(17, 6, 'Blue Airlines', 356);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(18, 6, 'Green Airlines', 306);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(19, 7, 'Red Airlines', 1007);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(20, 7, 'Blue Airlines', 357);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(21, 7, 'Green Airlines', 307);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(22, 8, 'Red Airlines', 1008);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(23, 8, 'Blue Airlines', 358);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(24, 8, 'Green Airlines', 308);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(25, 9, 'Red Airlines', 1009);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(26, 9, 'Blue Airlines', 359);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(27, 9, 'Green Airlines', 309);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(28, 10, 'Red Airlines', 1010);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(29, 10, 'Blue Airlines', 360);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(30, 10, 'Green Airlines', 310);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(31, 11, 'Red Airlines', 1011);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(32, 11, 'Blue Airlines', 361);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(33, 11, 'Green Airlines', 311);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(34, 12, 'Red Airlines', 1012);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(35, 12, 'Blue Airlines', 362);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(36, 12, 'Green Airlines', 312);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(37, 13, 'Red Airlines', 1013);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(38, 13, 'Blue Airlines', 363);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(39, 13, 'Green Airlines', 313);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(40, 14, 'Red Airlines', 1014);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(41, 14, 'Blue Airlines', 364);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(42, 14, 'Green Airlines', 314);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(43, 15, 'Red Airlines', 1015);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(44, 15, 'Blue Airlines', 365);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(45, 15, 'Green Airlines', 315);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(46, 16, 'Red Airlines', 1016);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(47, 16, 'Blue Airlines', 366);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(48, 16, 'Green Airlines', 316);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(49, 17, 'Red Airlines', 1017);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(50, 17, 'Blue Airlines', 367);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(51, 17, 'Green Airlines', 317);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(52, 18, 'Red Airlines', 1018);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(53, 18, 'Blue Airlines', 368);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(54, 18, 'Green Airlines', 318);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(55, 19, 'Red Airlines', 1019);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(56, 19, 'Blue Airlines', 369);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(57, 19, 'Green Airlines', 319);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(58, 20, 'Red Airlines', 1020);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(59, 20, 'Blue Airlines', 370);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(60, 20, 'Green Airlines', 320);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(61, 21, 'Red Airlines', 1021);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(62, 21, 'Blue Airlines', 371);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(63, 21, 'Green Airlines', 321);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(64, 22, 'Red Airlines', 1022);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(65, 22, 'Blue Airlines', 372);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(66, 22, 'Green Airlines', 322);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(67, 23, 'Red Airlines', 1023);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(68, 23, 'Blue Airlines', 373);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(69, 23, 'Green Airlines', 323);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(70, 24, 'Red Airlines', 1024);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(71, 24, 'Blue Airlines', 374);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(72, 24, 'Green Airlines', 324);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(73, 25, 'Red Airlines', 1025);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(74, 25, 'Blue Airlines', 375);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(75, 25, 'Green Airlines', 325);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(76, 26, 'Red Airlines', 1026);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(77, 26, 'Blue Airlines', 376);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(78, 26, 'Green Airlines', 326);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(79, 27, 'Red Airlines', 1027);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(80, 27, 'Blue Airlines', 377);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(81, 27, 'Green Airlines', 327);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(82, 28, 'Red Airlines', 1028);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(83, 28, 'Blue Airlines', 378);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(84, 28, 'Green Airlines', 328);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(85, 29, 'Red Airlines', 1029);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(86, 29, 'Blue Airlines', 379);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(87, 29, 'Green Airlines', 329);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(88, 30, 'Red Airlines', 1030);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(89, 30, 'Blue Airlines', 380);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(90, 30, 'Green Airlines', 330);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(91, 31, 'Red Airlines', 1031);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(92, 31, 'Blue Airlines', 381);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(93, 31, 'Green Airlines', 331);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(94, 32, 'Red Airlines', 1032);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(95, 32, 'Blue Airlines', 382);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(96, 32, 'Green Airlines', 332);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(97, 33, 'Red Airlines', 1033);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(98, 33, 'Blue Airlines', 383);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(99, 33, 'Green Airlines', 333);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(100, 34, 'Red Airlines', 1034);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(101, 34, 'Blue Airlines', 384);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(102, 34, 'Green Airlines', 334);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(103, 35, 'Red Airlines', 1035);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(104, 35, 'Blue Airlines', 385);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(105, 35, 'Green Airlines', 335);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(106, 36, 'Red Airlines', 1036);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(107, 36, 'Blue Airlines', 386);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(108, 36, 'Green Airlines', 336);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(109, 37, 'Red Airlines', 1037);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(110, 37, 'Blue Airlines', 387);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(111, 37, 'Green Airlines', 337);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(112, 38, 'Red Airlines', 1038);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(113, 38, 'Blue Airlines', 388);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(114, 38, 'Green Airlines', 338);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(115, 39, 'Red Airlines', 1039);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(116, 39, 'Blue Airlines', 389);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(117, 39, 'Green Airlines', 339);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(118, 40, 'Red Airlines', 1040);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(119, 40, 'Blue Airlines', 390);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(120, 40, 'Green Airlines', 340);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(121, 41, 'Red Airlines', 1041);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(122, 41, 'Blue Airlines', 391);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(123, 41, 'Green Airlines', 341);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(124, 42, 'Red Airlines', 1042);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(125, 42, 'Blue Airlines', 392);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(126, 42, 'Green Airlines', 342);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(127, 43, 'Red Airlines', 1043);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(128, 43, 'Blue Airlines', 393);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(129, 43, 'Green Airlines', 343);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(130, 44, 'Red Airlines', 1044);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(131, 44, 'Blue Airlines', 394);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(132, 44, 'Green Airlines', 344);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(133, 45, 'Red Airlines', 1045);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(134, 45, 'Blue Airlines', 395);
INSERT INTO flights (flightId, cityId, airline, price) VALUES(135, 45, 'Green Airlines', 345);

INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(1, 1, 'Grand Hotel Amsterdam', 505);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(2, 1, 'Little Amsterdam Hotel', 82);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(3, 2, 'Grand Hotel Andorra', 510);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(4, 2, 'Little Andorra Hotel', 84);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(5, 3, 'Grand Hotel Athens', 515);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(6, 3, 'Little Athens Hotel', 86);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(7, 4, 'Grand Hotel Belgrade', 520);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(8, 4, 'Little Belgrade Hotel', 88);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(9, 5, 'Grand Hotel Berlin', 525);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(10, 5, 'Little Berlin Hotel', 90);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(11, 6, 'Grand Hotel Bern', 530);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(12, 6, 'Little Bern Hotel', 92);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(13, 7, 'Grand Hotel Bratislava', 535);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(14, 7, 'Little Bratislava Hotel', 94);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(15, 8, 'Grand Hotel Brussels', 540);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(16, 8, 'Little Brussels Hotel', 96);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(17, 9, 'Grand Hotel Bucharest', 545);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(18, 9, 'Little Bucharest Hotel', 98);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(19, 10, 'Grand Hotel Budapest', 550);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(20, 10, 'Little Budapest Hotel', 100);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(21, 11, 'Grand Hotel Chisinau', 555);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(22, 11, 'Little Chisinau Hotel', 102);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(23, 12, 'Grand Hotel Copenhagen', 560);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(24, 12, 'Little Copenhagen Hotel', 104);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(25, 13, 'Grand Hotel Dublin', 565);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(26, 13, 'Little Dublin Hotel', 106);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(27, 14, 'Grand Hotel Helsinki', 570);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(28, 14, 'Little Helsinki Hotel', 108);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(29, 15, 'Grand Hotel Kiev', 575);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(30, 15, 'Little Kiev Hotel', 110);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(31, 16, 'Grand Hotel Lisbon', 580);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(32, 16, 'Little Lisbon Hotel', 112);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(33, 17, 'Grand Hotel Ljubljana', 585);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(34, 17, 'Little Ljubljana Hotel', 114);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(35, 18, 'Grand Hotel London', 590);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(36, 18, 'Little London Hotel', 116);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(37, 19, 'Grand Hotel Luxembourg', 595);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(38, 19, 'Little Luxembourg Hotel', 118);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(39, 20, 'Grand Hotel Madrid', 600);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(40, 20, 'Little Madrid Hotel', 120);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(41, 21, 'Grand Hotel Minsk', 605);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(42, 21, 'Little Minsk Hotel', 122);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(43, 22, 'Grand Hotel Monaco', 610);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(44, 22, 'Little Monaco Hotel', 124);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(45, 23, 'Grand Hotel Moscow', 615);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(46, 23, 'Little Moscow Hotel', 126);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(47, 24, 'Grand Hotel Nicosia', 620);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(48, 24, 'Little Nicosia Hotel', 128);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(49, 25, 'Grand Hotel Nuuk', 625);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(50, 25, 'Little Nuuk Hotel', 130);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(51, 26, 'Grand Hotel Oslo', 630);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(52, 26, 'Little Oslo Hotel', 132);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(53, 27, 'Grand Hotel Paris', 635);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(54, 27, 'Little Paris Hotel', 134);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(55, 28, 'Grand Hotel Podgorica', 640);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(56, 28, 'Little Podgorica Hotel', 136);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(57, 29, 'Grand Hotel Prague', 645);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(58, 29, 'Little Prague Hotel', 138);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(59, 30, 'Grand Hotel Reykjavik', 650);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(60, 30, 'Little Reykjavik Hotel', 140);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(61, 31, 'Grand Hotel Riga', 655);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(62, 31, 'Little Riga Hotel', 142);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(63, 32, 'Grand Hotel Rome', 660);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(64, 32, 'Little Rome Hotel', 144);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(65, 33, 'Grand Hotel San Marino', 665);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(66, 33, 'Little San Marino Hotel', 146);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(67, 34, 'Grand Hotel Sarajevo', 670);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(68, 34, 'Little Sarajevo Hotel', 148);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(69, 35, 'Grand Hotel Skopje', 675);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(70, 35, 'Little Skopje Hotel', 150);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(71, 36, 'Grand Hotel Sofia', 680);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(72, 36, 'Little Sofia Hotel', 152);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(73, 37, 'Grand Hotel Stockholm', 685);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(74, 37, 'Little Stockholm Hotel', 154);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(75, 38, 'Grand Hotel Tallinn', 690);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(76, 38, 'Little Tallinn Hotel', 156);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(77, 39, 'Grand Hotel Tirana', 695);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(78, 39, 'Little Tirana Hotel', 158);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(79, 40, 'Grand Hotel Vaduz', 700);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(80, 40, 'Little Vaduz Hotel', 160);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(81, 41, 'Grand Hotel Valletta', 705);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(82, 41, 'Little Valletta Hotel', 162);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(83, 42, 'Grand Hotel Vienna', 710);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(84, 42, 'Little Vienna Hotel', 164);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(85, 43, 'Grand Hotel Vilnius', 715);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(86, 43, 'Little Vilnius Hotel', 166);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(87, 44, 'Grand Hotel Warsaw', 720);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(88, 44, 'Little Warsaw Hotel', 168);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(89, 45, 'Grand Hotel Zagreb', 725);
INSERT INTO hotels (hotelId, cityId, hotel, price) VALUES(90, 45, 'Little Zagreb Hotel', 170);

INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(1, 1, 'Yellow Insurances', 308);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(2, 1, 'Blue Insurances', 57);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(3, 2, 'Yellow Insurances', 309);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(4, 2, 'Blue Insurances', 58);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(5, 3, 'Yellow Insurances', 310);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(6, 3, 'Blue Insurances', 59);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(7, 4, 'Yellow Insurances', 311);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(8, 4, 'Blue Insurances', 60);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(9, 5, 'Yellow Insurances', 312);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(10, 5, 'Blue Insurances', 61);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(11, 6, 'Yellow Insurances', 313);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(12, 6, 'Blue Insurances', 62);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(13, 7, 'Yellow Insurances', 314);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(14, 7, 'Blue Insurances', 63);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(15, 8, 'Yellow Insurances', 315);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(16, 8, 'Blue Insurances', 64);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(17, 9, 'Yellow Insurances', 316);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(18, 9, 'Blue Insurances', 65);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(19, 10, 'Yellow Insurances', 317);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(20, 10, 'Blue Insurances', 66);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(21, 11, 'Yellow Insurances', 318);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(22, 11, 'Blue Insurances', 67);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(23, 12, 'Yellow Insurances', 319);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(24, 12, 'Blue Insurances', 68);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(25, 13, 'Yellow Insurances', 320);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(26, 13, 'Blue Insurances', 69);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(27, 14, 'Yellow Insurances', 321);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(28, 14, 'Blue Insurances', 70);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(29, 15, 'Yellow Insurances', 322);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(30, 15, 'Blue Insurances', 71);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(31, 16, 'Yellow Insurances', 323);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(32, 16, 'Blue Insurances', 72);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(33, 17, 'Yellow Insurances', 324);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(34, 17, 'Blue Insurances', 73);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(35, 18, 'Yellow Insurances', 325);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(36, 18, 'Blue Insurances', 74);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(37, 19, 'Yellow Insurances', 326);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(38, 19, 'Blue Insurances', 75);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(39, 20, 'Yellow Insurances', 327);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(40, 20, 'Blue Insurances', 76);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(41, 21, 'Yellow Insurances', 328);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(42, 21, 'Blue Insurances', 77);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(43, 22, 'Yellow Insurances', 329);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(44, 22, 'Blue Insurances', 78);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(45, 23, 'Yellow Insurances', 330);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(46, 23, 'Blue Insurances', 79);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(47, 24, 'Yellow Insurances', 331);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(48, 24, 'Blue Insurances', 80);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(49, 25, 'Yellow Insurances', 332);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(50, 25, 'Blue Insurances', 81);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(51, 26, 'Yellow Insurances', 333);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(52, 26, 'Blue Insurances', 82);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(53, 27, 'Yellow Insurances', 334);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(54, 27, 'Blue Insurances', 83);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(55, 28, 'Yellow Insurances', 335);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(56, 28, 'Blue Insurances', 84);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(57, 29, 'Yellow Insurances', 336);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(58, 29, 'Blue Insurances', 85);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(59, 30, 'Yellow Insurances', 337);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(60, 30, 'Blue Insurances', 86);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(61, 31, 'Yellow Insurances', 338);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(62, 31, 'Blue Insurances', 87);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(63, 32, 'Yellow Insurances', 339);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(64, 32, 'Blue Insurances', 88);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(65, 33, 'Yellow Insurances', 340);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(66, 33, 'Blue Insurances', 89);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(67, 34, 'Yellow Insurances', 341);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(68, 34, 'Blue Insurances', 90);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(69, 35, 'Yellow Insurances', 342);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(70, 35, 'Blue Insurances', 91);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(71, 36, 'Yellow Insurances', 343);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(72, 36, 'Blue Insurances', 92);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(73, 37, 'Yellow Insurances', 344);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(74, 37, 'Blue Insurances', 93);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(75, 38, 'Yellow Insurances', 345);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(76, 38, 'Blue Insurances', 94);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(77, 39, 'Yellow Insurances', 346);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(78, 39, 'Blue Insurances', 95);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(79, 40, 'Yellow Insurances', 347);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(80, 40, 'Blue Insurances', 96);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(81, 41, 'Yellow Insurances', 348);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(82, 41, 'Blue Insurances', 97);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(83, 42, 'Yellow Insurances', 349);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(84, 42, 'Blue Insurances', 98);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(85, 43, 'Yellow Insurances', 350);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(86, 43, 'Blue Insurances', 99);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(87, 44, 'Yellow Insurances', 351);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(88, 44, 'Blue Insurances', 100);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(89, 45, 'Yellow Insurances', 352);
INSERT INTO insurances (insuranceId, cityId, company, price) VALUES(90, 45, 'Blue Insurances', 101);
