import sys
from PyQt6.QtWidgets import (QApplication, QMainWindow, QTableView, 
                             QVBoxLayout, QHBoxLayout, QWidget, QLabel, 
                             QMessageBox, QSplitter, QComboBox)
from PyQt6.QtSql import QSqlDatabase, QSqlQueryModel, QSqlQuery
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QFont

# MySQL
MYSQL_HOST = "127.0.0.1"
MYSQL_DB = "insurance"
MYSQL_USER = "root"
MYSQL_PASS = "1234"

# PostgreSQL
PG_HOST = "127.0.0.1"
PG_DB = "insurance"
PG_USER = "postgres"
PG_PASS = "1234"

# SQL Запросы (настраиваемые)
SQL_EMPLOYEES = "SELECT employee_id, full_name, passport, position, count_policyholders FROM employees"

# Запрос для подготовленного Statement 
SQL_CLAIMS_SAFE = """
    SELECT c.claim_id, p.policy_number, c.description, p.full_name, 
           p.birth_date, p.passport, c.event_date, c.payout 
    FROM claims c 
    JOIN policyholders p ON c.policy_number = p.policy_number 
    WHERE p.employee_id = :emp_id
"""

# Запрос для прямого форматирования
SQL_CLAIMS_RAW = """
    SELECT c.claim_id, p.policy_number, c.description, p.full_name, 
           p.birth_date, p.passport, c.event_date, c.payout 
    FROM claims c 
    JOIN policyholders p ON c.policy_number = p.policy_number 
    WHERE p.employee_id = {0}
"""

class DatabaseSection(QWidget):
    def __init__(self, title, db_connection_name):
        super().__init__()
        self.db_name = db_connection_name
        self.init_ui(title)
        self.load_employees()

    def init_ui(self, title):
        layout = QVBoxLayout()
        layout.setContentsMargins(10, 10, 10, 10)

        # Шрифты
        font_title = QFont("Bahnschrift", 24)
        font_subtitle = QFont("Bahnschrift", 14)
        font_labels = QFont("Bahnschrift", 12)

        # Заголовок
        lbl_title = QLabel(title)
        lbl_title.setFont(font_title)
        layout.addWidget(lbl_title)

        # Подзаголовок 1
        lbl_sub1 = QLabel("1: Выберите сотрудника (ComboBox - аналог C#):")
        lbl_sub1.setFont(font_subtitle)
        layout.addWidget(lbl_sub1)

        # --- ЗАМЕНА ТАБЛИЦЫ НА COMBOBOX ---
        self.employee_combo = QComboBox()
        self.employee_combo.setFont(font_labels)
        self.employee_combo.setMinimumHeight(35)
        # Сигнал изменения выбора (аналог SelectedIndexChanged)
        self.employee_combo.currentIndexChanged.connect(self.on_employee_selected)
        layout.addWidget(self.employee_combo)
        # ----------------------------------

        # Подзаголовок 2
        lbl_sub2 = QLabel("2: Таблица выделяет Claims Policyholders и выводит случившиеся происшествия")
        lbl_sub2.setFont(font_subtitle)
        layout.addWidget(lbl_sub2)

        # Метки над нижними таблицами
        labels_layout = QHBoxLayout()
        lbl_linq = QLabel("Запрос через Bind (Аналог LINQ)")
        lbl_linq.setFont(font_labels)
        lbl_linq.setStyleSheet("color: gray;")
        lbl_linq.setAlignment(Qt.AlignmentFlag.AlignCenter)
        
        lbl_sql = QLabel("Прямой SQL запрос (Аналог FromSqlRaw)")
        lbl_sql.setFont(font_labels)
        lbl_sql.setStyleSheet("color: gray;")
        lbl_sql.setAlignment(Qt.AlignmentFlag.AlignCenter)

        labels_layout.addWidget(lbl_linq)
        labels_layout.addWidget(lbl_sql)
        layout.addLayout(labels_layout)

        # Нижние таблицы (Detail)
        tables_layout = QHBoxLayout()
        self.detail_view_safe = QTableView()
        self.detail_view_safe.setSelectionBehavior(QTableView.SelectionBehavior.SelectRows)
        self.detail_view_raw = QTableView()
        self.detail_view_raw.setSelectionBehavior(QTableView.SelectionBehavior.SelectRows)

        tables_layout.addWidget(self.detail_view_safe)
        tables_layout.addWidget(self.detail_view_raw)
        layout.addLayout(tables_layout)

        self.setLayout(layout)

        # Модели
        self.master_model = QSqlQueryModel()
        self.detail_safe_model = QSqlQueryModel()
        self.detail_raw_model = QSqlQueryModel()

        # Привязываем модель к комбобоксу
        self.employee_combo.setModel(self.master_model)
        self.detail_view_safe.setModel(self.detail_safe_model)
        self.detail_view_raw.setModel(self.detail_raw_model)

    def load_employees(self):
        """ Загрузка списка сотрудников в ComboBox """
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.exec(SQL_EMPLOYEES)
        
        self.master_model.setQuery(query)
        
        # Указываем комбобоксу, какую колонку показывать пользователю (Full Name - индекс 1)
        # employee_id (индекс 0) будет "спрятан" за этим текстом
        self.employee_combo.setModelColumn(1) 

    def on_employee_selected(self, index):
        """ Обработка выбора в ComboBox (аналог SelectedValue) """
        if index < 0:
            return

        # Получаем ID сотрудника из скрытой колонки 0 той же строки
        employee_id = self.master_model.record(index).value("employee_id")
        
        if employee_id is None:
            return

        db = QSqlDatabase.database(self.db_name)

        # Безопасный запрос (Bind)
        query_safe = QSqlQuery(db)
        query_safe.prepare(SQL_CLAIMS_SAFE)
        query_safe.bindValue(":emp_id", employee_id)
        query_safe.exec()
        self.detail_safe_model.setQuery(query_safe)
        self._set_claims_headers(self.detail_safe_model)

        # Прямой SQL запрос
        query_raw = QSqlQuery(db)
        formatted_sql = SQL_CLAIMS_RAW.format(employee_id)
        query_raw.exec(formatted_sql)
        self.detail_raw_model.setQuery(query_raw)
        self._set_claims_headers(self.detail_raw_model)

    def _set_claims_headers(self, model):
        headers = ["ID", "ID Потерпевшего", "Описание", "ФИО", "День Рождения", "Паспорт", "Дата происшествия", "Выплата"]
        for i, header in enumerate(headers):
            model.setHeaderData(i, Qt.Orientation.Horizontal, header)


#  ГЛАВНОЕ ОКНО ПРИЛОЖЕНИЯ
class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Таблички (PyQt6)")
        self.resize(1200, 800)

        # Инициализация БД
        if not self.init_databases():
            sys.exit(1)

        # Основной Layout через QSplitter (чтобы можно было менять размер панелей)
        splitter = QSplitter(Qt.Orientation.Vertical)

        # Создаем две одинаковые секции для разных БД
        self.mysql_section = DatabaseSection("MySQL", "mysql_conn")
        self.pg_section = DatabaseSection("PostgreSQL", "pg_conn")

        splitter.addWidget(self.mysql_section)
        splitter.addWidget(self.pg_section)

        self.setCentralWidget(splitter)

    def init_databases(self):
        """ Подключение к обеим базам данных """
        # Подключение MySQL
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn") # или QMYSQL
        db_ms.setHostName(MYSQL_HOST)
        db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER)
        db_ms.setPassword(MYSQL_PASS)

        if not db_ms.open():
            QMessageBox.critical(self, "Ошибка MySQL", f"Не удалось подключиться к MySQL:\n{db_ms.lastError().text()}")
            return False

        # Подключение PostgreSQL
        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST)
        db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER)
        db_pg.setPassword(PG_PASS)

        if not db_pg.open():
            QMessageBox.critical(self, "Ошибка PostgreSQL", f"Не удалось подключиться к PostgreSQL:\n{db_pg.lastError().text()}")
            return False

        return True

if __name__ == "__main__":
    app = QApplication(sys.argv)
    
    app.setFont(QFont("Times New Roman", 12))
    
    window = MainWindow()
    window.showMaximized()
    sys.exit(app.exec())