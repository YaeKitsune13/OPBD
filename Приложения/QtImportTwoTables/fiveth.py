import sys
from PyQt6.QtWidgets import (QApplication, QMainWindow, QVBoxLayout, QHBoxLayout,
                             QWidget, QLabel, QMessageBox, QSplitter,
                             QComboBox, QScrollArea, QFrame)
from PyQt6.QtSql import QSqlDatabase, QSqlQueryModel, QSqlQuery
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QFont

# Настройки подключения
MYSQL_HOST, MYSQL_DB, MYSQL_USER, MYSQL_PASS = "127.0.0.1", "insurance", "root", "1234"
PG_HOST, PG_DB, PG_USER, PG_PASS = "127.0.0.1", "insurance", "postgres", "1234"

# --- ОБНОВЛЕННЫЕ SQL ЗАПРОСЫ ---
# Используем UNION для добавления варианта "Все"
SQL_GET_EMPLOYEES = """
    SELECT 0 AS employee_id, 'Все' AS full_name
    UNION ALL
    SELECT employee_id, full_name FROM employees
"""

# Для клиентов: добавляем "Все", фильтруем по сотруднику если emp_id != 0
SQL_GET_POLICYHOLDERS_ALL = """
    SELECT '0' AS policy_number, 'Все' AS full_name
    UNION ALL
    SELECT policy_number, full_name FROM policyholders
"""
SQL_GET_POLICYHOLDERS_FILTER = """
    SELECT '0' AS policy_number, 'Все' AS full_name
    UNION ALL
    SELECT policy_number, full_name FROM policyholders WHERE employee_id = :emp_id
"""

# Для страховых случаев: разные сценарии выборки
SQL_CLAIMS_EVERYTHING = "SELECT claim_id, description, event_date, payout, policy_number FROM claims"
SQL_CLAIMS_BY_EMPLOYEE = """
    SELECT c.claim_id, c.description, c.event_date, c.payout, c.policy_number 
    FROM claims c
    JOIN policyholders p ON c.policy_number = p.policy_number
    WHERE p.employee_id = :emp_id
"""
SQL_CLAIMS_BY_POLICY = "SELECT claim_id, description, event_date, payout, policy_number FROM claims WHERE policy_number = :policy_num"


class ClaimCard(QFrame):
    def __init__(self, policy_num, claim_id, description, date, payout):
        super().__init__()
        self.setFrameShape(QFrame.Shape.NoFrame)
        self.setFixedHeight(100)
        self.setContentsMargins(5, 5, 5, 5)
        self.setStyleSheet("""
            ClaimCard {
                background-color: #F9F9F9;
                border: 1px solid #CCCCCC;
                border-radius: 8px;
            }
        """)

        outer = QHBoxLayout(self)
        outer.setContentsMargins(10, 10, 10, 10)
        outer.setSpacing(0)

        placeholder = QLabel("📄")
        placeholder.setFixedSize(70, 70)
        placeholder.setAlignment(Qt.AlignmentFlag.AlignCenter)
        placeholder.setStyleSheet("background-color: #CCCCCC; border-radius: 8px; font-size: 28px;")
        outer.addWidget(placeholder, 0, Qt.AlignmentFlag.AlignVCenter)

        info_layout = QVBoxLayout()
        info_layout.setContentsMargins(15, 0, 15, 0)
        info_layout.setSpacing(4)

        lbl_desc = QLabel(str(description))
        lbl_desc.setFont(QFont("Segoe UI", 10, QFont.Weight.Bold))
        lbl_desc.setStyleSheet("color: #1a1a1a; font-size: 13px;")

        # Форматирование даты
        try:
            if hasattr(date, "toPyDate"): date_str = date.toPyDate().strftime("%d.%m.%Y")
            elif hasattr(date, "toString"): date_str = date.toString("dd.MM.yyyy")
            else: date_str = str(date)
        except: date_str = str(date)

        lbl_date = QLabel(f"Дата: {date_str}")
        lbl_date.setStyleSheet("color: #888888; font-size: 11px;")

        lbl_policy = QLabel(f"Полис: {policy_num}")
        lbl_policy.setStyleSheet("color: #888888; font-size: 11px;")

        info_layout.addWidget(lbl_desc)
        info_layout.addWidget(lbl_date)
        info_layout.addWidget(lbl_policy)
        outer.addLayout(info_layout, 1)

        payout_badge = QFrame()
        payout_badge.setFixedSize(110, 70)
        payout_badge.setStyleSheet("background-color: #4CAF50; border-radius: 8px;")
        badge_layout = QVBoxLayout(payout_badge)
        
        lbl_payout_title = QLabel("Выплата")
        lbl_payout_title.setAlignment(Qt.AlignmentFlag.AlignCenter)
        lbl_payout_title.setStyleSheet("color: white; font-size: 11px; background: transparent;")

        try: payout_formatted = f"{float(payout):,.0f}".replace(",", " ")
        except: payout_formatted = str(payout)

        lbl_payout_value = QLabel(f"{payout_formatted} ₽")
        lbl_payout_value.setAlignment(Qt.AlignmentFlag.AlignCenter)
        lbl_payout_value.setStyleSheet("color: white; font-size: 14px; font-weight: bold; background: transparent;")

        badge_layout.addWidget(lbl_payout_title)
        badge_layout.addWidget(lbl_payout_value)
        outer.addWidget(payout_badge, 0, Qt.AlignmentFlag.AlignVCenter)


class DatabaseSection(QWidget):
    def __init__(self, title, db_connection_name):
        super().__init__()
        self.db_name = db_connection_name
        self.current_emp_id = 0
        self.init_ui(title)
        self.load_employees()

    def init_ui(self, title):
        layout = QVBoxLayout()
        lbl_title = QLabel(title)
        lbl_title.setFont(QFont("Bahnschrift", 22))
        lbl_title.setStyleSheet("color: #2980b9;")
        layout.addWidget(lbl_title)

        grid_layout = QHBoxLayout()
        
        v_emp = QVBoxLayout()
        v_emp.addWidget(QLabel("1. Сотрудник:"))
        self.combo_emp = QComboBox()
        v_emp.addWidget(self.combo_emp)

        v_hold = QVBoxLayout()
        v_hold.addWidget(QLabel("2. Клиент:"))
        self.combo_holder = QComboBox()
        v_hold.addWidget(self.combo_holder)

        grid_layout.addLayout(v_emp)
        grid_layout.addLayout(v_hold)
        layout.addLayout(grid_layout)

        self.scroll = QScrollArea()
        self.scroll.setWidgetResizable(True)
        self.scroll_content = QWidget()
        self.cards_layout = QVBoxLayout(self.scroll_content)
        self.cards_layout.setAlignment(Qt.AlignmentFlag.AlignTop)
        self.scroll.setWidget(self.scroll_content)
        layout.addWidget(self.scroll)

        self.setLayout(layout)

        self.model_emp = QSqlQueryModel()
        self.model_holder = QSqlQueryModel()
        self.combo_emp.setModel(self.model_emp)
        self.combo_holder.setModel(self.model_holder)

        self.combo_emp.currentIndexChanged.connect(self.on_employee_changed)
        self.combo_holder.currentIndexChanged.connect(self.on_policyholder_changed)

    def load_employees(self):
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.exec(SQL_GET_EMPLOYEES)
        self.model_emp.setQuery(query)
        self.combo_emp.setModelColumn(1)

    def on_employee_changed(self, index):
        if index < 0: return
        self.current_emp_id = self.model_emp.record(index).value("employee_id")

        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        
        if self.current_emp_id == 0:
            query.exec(SQL_GET_POLICYHOLDERS_ALL)
        else:
            query.prepare(SQL_GET_POLICYHOLDERS_FILTER)
            query.bindValue(":emp_id", self.current_emp_id)
            query.exec()

        self.model_holder.setQuery(query)
        self.combo_holder.setModelColumn(1)
        # Принудительно сбрасываем индекс, чтобы сработало событие "Все"
        self.combo_holder.setCurrentIndex(0)

    def on_policyholder_changed(self, index):
        # Очистка карточек
        while self.cards_layout.count():
            child = self.cards_layout.takeAt(0)
            if child.widget(): child.widget().deleteLater()

        if index < 0: return

        policy_num = str(self.model_holder.record(index).value("policy_number"))
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)

        # ЛОГИКА ВЫБОРКИ CLAIMS:
        if policy_num == "0":  # Выбрано "Все" в клиентах
            if self.current_emp_id == 0: # И "Все" в сотрудниках
                query.prepare(SQL_CLAIMS_EVERYTHING)
            else: # "Все" клиенты конкретного сотрудника
                query.prepare(SQL_CLAIMS_BY_EMPLOYEE)
                query.bindValue(":emp_id", self.current_emp_id)
        else: # Конкретный клиент
            query.prepare(SQL_CLAIMS_BY_POLICY)
            query.bindValue(":policy_num", policy_num)
        
        query.exec()

        while query.next():
            card = ClaimCard(
                policy_num=query.value("policy_number"),
                claim_id=query.value("claim_id"),
                description=query.value("description"),
                date=query.value("event_date"),
                payout=query.value("payout")
            )
            self.cards_layout.addWidget(card)


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Insurance Cards (Filter: All)")
        self.resize(1200, 800)
        if not self.init_databases(): sys.exit(1)

        splitter = QSplitter(Qt.Orientation.Horizontal)
        splitter.addWidget(DatabaseSection("MySQL", "mysql_conn"))
        splitter.addWidget(DatabaseSection("PostgreSQL", "pg_conn"))
        self.setCentralWidget(splitter)

    def init_databases(self):
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn") # или QMYSQL
        db_ms.setHostName(MYSQL_HOST); db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER); db_ms.setPassword(MYSQL_PASS)

        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST); db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER); db_pg.setPassword(PG_PASS)

        if not db_ms.open() or not db_pg.open():
            QMessageBox.critical(self, "Error", "Database Connection Failed")
            return False
        return True

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setFont(QFont("Segoe UI", 10))
    window = MainWindow()
    window.show()
    sys.exit(app.exec())