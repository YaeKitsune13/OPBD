SELECT COUNT(DISTINCT j_no) AS total_projects
FROM SPJ
WHERE s_no = 'S1';

SELECT DISTINCT j_no
FROM SPJ
WHERE s_no IN (
    SELECT DISTINCT SPJ.s_no
    FROM SPJ
             JOIN P ON SPJ.p_no = P.p_no
    WHERE P.color = 'Красный'
);

SELECT DISTINCT P.*
FROM P
         JOIN SPJ ON P.p_no = SPJ.p_no
         JOIN J ON SPJ.j_no = J.j_no
WHERE J.city = 'Лондон';

SELECT DISTINCT SPJ.p_no
FROM SPJ
         JOIN S ON SPJ.s_no = S.s_no
WHERE S.city = 'Лондон';

SELECT p_no
FROM P
WHERE NOT EXISTS (
    SELECT j_no
    FROM J
    WHERE city = 'Лондон'
      AND NOT EXISTS (
        SELECT *
        FROM SPJ
        WHERE SPJ.p_no = P.p_no
          AND SPJ.j_no = J.j_no
    )
);

SELECT DISTINCT j_no
FROM SPJ
WHERE s_no = 'S1';

SELECT DISTINCT P.p_no, P.p_name
FROM P
         JOIN SPJ ON P.p_no = SPJ.p_no
         JOIN J ON SPJ.j_no = J.j_no
WHERE J.city = 'Лондон';

SELECT j_no, j_name
FROM J
WHERE city = (
    SELECT MIN(city)
    FROM J
);

SELECT DISTINCT P.color
FROM P
         JOIN SPJ ON P.p_no = SPJ.p_no
WHERE SPJ.s_no = 'S1';

SELECT DISTINCT S.s_no, S.s_name
FROM S
         JOIN SPJ ON S.s_no = SPJ.s_no
WHERE SPJ.p_no = 'P1'
  AND SPJ.qty > (
    SELECT AVG(qty)
    FROM SPJ AS SPJ2
    WHERE SPJ2.p_no = 'P1'
      AND SPJ2.j_no = SPJ.j_no
);

SELECT DISTINCT J.j_name
FROM J
         JOIN SPJ ON J.j_no = SPJ.j_no
WHERE SPJ.s_no = 'S1';

SELECT DISTINCT SPJ.p_no
FROM SPJ
         JOIN S ON SPJ.s_no = S.s_no
         JOIN J ON SPJ.j_no = J.j_no
WHERE S.city = J.city;