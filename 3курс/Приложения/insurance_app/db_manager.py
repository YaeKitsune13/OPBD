"""
Менеджер подключения к MySQL и базовые CRUD-операции.
Вся работа с БД сосредоточена здесь, чтобы UI ничего не знал про SQL.
"""
from __future__ import annotations

import pymysql
import pymysql.cursors


class DatabaseError(Exception):
    """Ошибка работы с базой данных в формате, понятном пользователю."""


class DatabaseManager:
    def __init__(self) -> None:
        self._conn: pymysql.connections.Connection | None = None

    @property
    def is_connected(self) -> bool:
        return self._conn is not None and self._conn.open

    def connect(self, host: str, port: int, user: str, password: str, database: str) -> None:
        try:
            self._conn = pymysql.connect(
                host=host,
                port=port,
                user=user,
                password=password,
                database=database,
                charset="utf8mb4",
                cursorclass=pymysql.cursors.DictCursor,
                autocommit=True,
            )
        except pymysql.MySQLError as exc:
            raise DatabaseError(f"Не удалось подключиться к базе данных: {exc}") from exc

    def close(self) -> None:
        if self._conn is not None:
            try:
                self._conn.close()
            except pymysql.MySQLError:
                pass
            self._conn = None

    def _cursor(self):
        if not self.is_connected:
            raise DatabaseError("Нет подключения к базе данных")
        return self._conn.cursor()

    # ---------- чтение ----------

    def fetch_all(self, table: str, order_by: str | None = None) -> list[dict]:
        query = f"SELECT * FROM `{table}`"
        if order_by:
            query += f" ORDER BY `{order_by}`"
        return self._run_select(query)

    def fetch_lookup(self, table: str, key_col: str, label_col: str) -> list[tuple]:
        """Список (ключ, подпись) для заполнения выпадающих списков (внешние ключи)."""
        rows = self._run_select(
            f"SELECT `{key_col}` AS k, `{label_col}` AS v FROM `{table}` ORDER BY `{label_col}`"
        )
        return [(r["k"], r["v"]) for r in rows]

    def run_query(self, query: str, params=None) -> list[dict]:
        return self._run_select(query, params)

    def _run_select(self, query: str, params=None) -> list[dict]:
        try:
            with self._cursor() as cur:
                cur.execute(query, params or ())
                return list(cur.fetchall())
        except pymysql.MySQLError as exc:
            raise DatabaseError(self._friendly_message(exc)) from exc

    # ---------- запись ----------

    def insert(self, table: str, values: dict) -> int | None:
        cols = ", ".join(f"`{c}`" for c in values)
        placeholders = ", ".join(["%s"] * len(values))
        query = f"INSERT INTO `{table}` ({cols}) VALUES ({placeholders})"
        try:
            with self._cursor() as cur:
                cur.execute(query, list(values.values()))
                return cur.lastrowid or None
        except pymysql.MySQLError as exc:
            raise DatabaseError(self._friendly_message(exc)) from exc

    def update(self, table: str, pk_col: str, pk_val, values: dict) -> None:
        if not values:
            return
        set_clause = ", ".join(f"`{c}` = %s" for c in values)
        query = f"UPDATE `{table}` SET {set_clause} WHERE `{pk_col}` = %s"
        try:
            with self._cursor() as cur:
                cur.execute(query, list(values.values()) + [pk_val])
        except pymysql.MySQLError as exc:
            raise DatabaseError(self._friendly_message(exc)) from exc

    def delete(self, table: str, pk_col: str, pk_val) -> None:
        query = f"DELETE FROM `{table}` WHERE `{pk_col}` = %s"
        try:
            with self._cursor() as cur:
                cur.execute(query, [pk_val])
        except pymysql.MySQLError as exc:
            raise DatabaseError(self._friendly_message(exc)) from exc

    @staticmethod
    def _friendly_message(exc: pymysql.MySQLError) -> str:
        if exc.args:
            code = exc.args[0]
            message = exc.args[1] if len(exc.args) > 1 else str(exc)
            if code == 1644:
                # SIGNAL SQLSTATE '45000' из пользовательского триггера
                return str(message)
            if code in (1451, 1452):
                return "Операция нарушает связь с другой таблицей (внешний ключ)."
            if code == 1062:
                return "Запись с таким ключом уже существует."
            return f"Ошибка базы данных: {message}"
        return f"Ошибка базы данных: {exc}"
