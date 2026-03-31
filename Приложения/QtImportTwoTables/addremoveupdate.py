import sys
from PyQt6.QtWidgets import (
    QApplication, QMainWindow, QTableView, QVBoxLayout, QHBoxLayout,
    QWidget, QLabel, QMessageBox, QSplitter, QHeaderView, QComboBox,
    QPushButton, QDialog, QFormLayout, QLineEdit, QDateEdit, QDoubleSpinBox,
    QDialogButtonBox
)
from PyQt6.QtSql import QSqlDatabase, QSqlQueryModel, QSqlQuery
from PyQt6.QtCore import Qt, QDate
from PyQt6.QtGui import QFont

# Настройки подключения
MYSQL_HOST, MYSQL_DB, MYSQL_USER, MYSQL_PASS = "127.0.0.1", "insurance", "root", "1234"
PG_HOST,    PG_DB,    PG_USER,    PG_PASS    = "127.0.0.1", "insurance", "postgres", "1234"

# SQL
SQL_GET_EMPLOYEES = """
    SELECT 0 AS employee_id, 'Все' AS full_name
    UNION ALL
    SELECT employee_id, full_name FROM employees
"""
SQL_GET_POLICYHOLDERS_ALL = """
    SELECT '0' AS policy_number, 'Все' AS full_name
    UNION ALL
    SELECT policy_number, full_name FROM policyholders
"""
SQL_GET_POLICYHOLDERS_BY_EMP = """
    SELECT '0' AS policy_number, 'Все' AS full_name
    UNION ALL
    SELECT policy_number, full_name FROM policyholders WHERE employee_id = :emp_id
"""
SQL_CLAIMS_ALL       = "SELECT claim_id, policy_number, description, event_date, payout FROM claims"
SQL_CLAIMS_BY_EMP    = """
    SELECT c.claim_id, c.policy_number, c.description, c.event_date, c.payout
    FROM claims c JOIN policyholders p ON c.policy_number = p.policy_number
    WHERE p.employee_id = :emp_id
"""
SQL_CLAIMS_BY_POLICY = "SELECT claim_id, policy_number, description, event_date, payout FROM claims WHERE policy_number = :policy_num"

CLAIMS_HEADERS = ["ID случая", "№ Полиса", "Описание", "Дата", "Выплата"]


class ClaimDialog(QDialog):
    def __init__(self, db_name, data=None, parent=None):
        super().__init__(parent)
        self.db_name = db_name
        self.setWindowTitle("Страховой случай")
        self.setMinimumWidth(400)

        form = QFormLayout(self)

        self.policy_combo = QComboBox()
        self._load_policies()
        form.addRow("№ Полиса:", self.policy_combo)

        self.desc_edit = QLineEdit()
        form.addRow("Описание:", self.desc_edit)

        self.date_edit = QDateEdit(QDate.currentDate())
        self.date_edit.setCalendarPopup(True)
        self.date_edit.setDisplayFormat("yyyy-MM-dd")
        form.addRow("Дата события:", self.date_edit)

        self.payout_spin = QDoubleSpinBox()
        self.payout_spin.setRange(0, 999_999_999)
        self.payout_spin.setDecimals(2)
        form.addRow("Выплата:", self.payout_spin)

        buttons = QDialogButtonBox(
            QDialogButtonBox.StandardButton.Ok | QDialogButtonBox.StandardButton.Cancel
        )
        buttons.accepted.connect(self.accept)
        buttons.rejected.connect(self.reject)
        form.addRow(buttons)

        if data:
            idx = self.policy_combo.findText(str(data[1]))
            if idx >= 0:
                self.policy_combo.setCurrentIndex(idx)
            self.desc_edit.setText(str(data[2]))
            self.date_edit.setDate(QDate.fromString(str(data[3]), "yyyy-MM-dd"))
            self.payout_spin.setValue(float(data[4]))

    def _load_policies(self):
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        q.exec("SELECT policy_number FROM policyholders")
        while q.next():
            self.policy_combo.addItem(str(q.value(0)))

    def get_values(self):
        return (
            self.policy_combo.currentText(),
            self.desc_edit.text(),
            self.date_edit.date().toString("yyyy-MM-dd"),
            self.payout_spin.value()
        )


def make_crud_buttons(add_cb, edit_cb, del_cb):
    row = QHBoxLayout()
    btn_add  = QPushButton("Добавить")
    btn_edit = QPushButton("Изменить")
    btn_del  = QPushButton("Удалить")

    for b in (btn_add, btn_edit, btn_del):
        b.setFixedHeight(32)

    btn_add.setStyleSheet("color: #27ae60; font-weight: bold;")
    btn_del.setStyleSheet("color: #c0392b; font-weight: bold;")

    btn_add.clicked.connect(add_cb)
    btn_edit.clicked.connect(edit_cb)
    btn_del.clicked.connect(del_cb)

    row.addWidget(btn_add)
    row.addWidget(btn_edit)
    row.addWidget(btn_del)
    row.addStretch()
    return row


class DatabaseSection(QWidget):
    def __init__(self, title, db_name):
        super().__init__()
        self.db_name = db_name
        self.current_emp_id = 0
        layout = QVBoxLayout(self)

        header = QLabel(title)
        header.setFont(QFont("Bahnschrift", 18, QFont.Weight.Bold))
        header.setAlignment(Qt.AlignmentFlag.AlignCenter)
        header.setStyleSheet("color: #2c3e50; padding: 6px;")
        layout.addWidget(header)

        font_label = QFont("Segoe UI", 10)

        filter_row = QHBoxLayout()
        filter_row.addWidget(QLabel("Сотрудник:", font=font_label))
        self.combo_employee = QComboBox()
        self.combo_employee.setMinimumWidth(180)
        filter_row.addWidget(self.combo_employee)
        filter_row.addWidget(QLabel("Клиент:", font=font_label))
        self.combo_holder = QComboBox()
        self.combo_holder.setMinimumWidth(180)
        filter_row.addWidget(self.combo_holder)
        filter_row.addStretch()
        layout.addLayout(filter_row)

        layout.addLayout(make_crud_buttons(self.add_record, self.edit_record, self.del_record))

        self.view = QTableView()
        self.view.horizontalHeader().setSectionResizeMode(QHeaderView.ResizeMode.Stretch)
        self.view.setSelectionBehavior(QTableView.SelectionBehavior.SelectRows)
        layout.addWidget(self.view)

        self.model_emp    = QSqlQueryModel()
        self.model_holder = QSqlQueryModel()
        self.model_claims = QSqlQueryModel()

        self.combo_employee.setModel(self.model_emp)
        self.combo_holder.setModel(self.model_holder)
        self.view.setModel(self.model_claims)

        self.combo_employee.currentIndexChanged.connect(self.on_employee_changed)
        self.combo_holder.currentIndexChanged.connect(self.on_holder_changed)

        self.load_employees()

    def load_employees(self):
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        q.exec(SQL_GET_EMPLOYEES)
        self.model_emp.setQuery(q)
        self.combo_employee.setModelColumn(1)

    def on_employee_changed(self, index):
        if index < 0:
            return
        self.current_emp_id = self.model_emp.record(index).value("employee_id")
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        if self.current_emp_id == 0:
            q.exec(SQL_GET_POLICYHOLDERS_ALL)
        else:
            q.prepare(SQL_GET_POLICYHOLDERS_BY_EMP)
            q.bindValue(":emp_id", self.current_emp_id)
            q.exec()
        self.model_holder.setQuery(q)
        self.combo_holder.setModelColumn(1)
        self.combo_holder.setCurrentIndex(0)

    def on_holder_changed(self, index):
        if index < 0:
            return
        policy_num = str(self.model_holder.record(index).value("policy_number"))
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        if policy_num == "0":
            if self.current_emp_id == 0:
                q.exec(SQL_CLAIMS_ALL)
            else:
                q.prepare(SQL_CLAIMS_BY_EMP)
                q.bindValue(":emp_id", self.current_emp_id)
                q.exec()
        else:
            q.prepare(SQL_CLAIMS_BY_POLICY)
            q.bindValue(":policy_num", policy_num)
            q.exec()
        self.model_claims.setQuery(q)
        for i, h in enumerate(CLAIMS_HEADERS):
            self.model_claims.setHeaderData(i, Qt.Orientation.Horizontal, h)

    def _selected_row_data(self):
        idx = self.view.currentIndex()
        if not idx.isValid():
            return None
        r = idx.row()
        return [self.model_claims.record(r).value(c) for c in range(self.model_claims.columnCount())]

    def _refresh(self):
        self.on_holder_changed(self.combo_holder.currentIndex())

    def add_record(self):
        dlg = ClaimDialog(self.db_name, parent=self)
        if dlg.exec() != QDialog.DialogCode.Accepted:
            return
        v = dlg.get_values()
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        q.prepare("INSERT INTO claims (policy_number, description, event_date, payout) VALUES (:pn,:d,:dt,:pay)")
        q.bindValue(":pn",  v[0])
        q.bindValue(":d",   v[1])
        q.bindValue(":dt",  v[2])
        q.bindValue(":pay", v[3])
        if not q.exec():
            QMessageBox.critical(self, "Ошибка", q.lastError().text())
        self._refresh()

    def edit_record(self):
        data = self._selected_row_data()
        if data is None:
            QMessageBox.information(self, "Внимание", "Выберите строку для редактирования.")
            return
        dlg = ClaimDialog(self.db_name, data=data, parent=self)
        if dlg.exec() != QDialog.DialogCode.Accepted:
            return
        v = dlg.get_values()
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        q.prepare("UPDATE claims SET policy_number=:pn, description=:d, event_date=:dt, payout=:pay WHERE claim_id=:id")
        q.bindValue(":pn",  v[0])
        q.bindValue(":d",   v[1])
        q.bindValue(":dt",  v[2])
        q.bindValue(":pay", v[3])
        q.bindValue(":id",  data[0])
        if not q.exec():
            QMessageBox.critical(self, "Ошибка", q.lastError().text())
        self._refresh()

    def del_record(self):
        data = self._selected_row_data()
        if data is None:
            QMessageBox.information(self, "Внимание", "Выберите строку для удаления.")
            return
        if QMessageBox.question(
            self, "Подтверждение",
            f"Удалить случай №{data[0]}?",
            QMessageBox.StandardButton.Yes | QMessageBox.StandardButton.No
        ) != QMessageBox.StandardButton.Yes:
            return
        db = QSqlDatabase.database(self.db_name)
        q = QSqlQuery(db)
        q.prepare("DELETE FROM claims WHERE claim_id = :id")
        q.bindValue(":id", data[0])
        if not q.exec():
            QMessageBox.critical(self, "Ошибка", q.lastError().text())
        self._refresh()


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Insurance Database Management — CRUD")
        self.resize(1300, 750)
        if not self.init_databases():
            sys.exit(1)

        splitter = QSplitter(Qt.Orientation.Horizontal)
        splitter.addWidget(DatabaseSection("MySQL", "mysql_conn"))
        splitter.addWidget(DatabaseSection("PostgreSQL", "pg_conn"))
        self.setCentralWidget(splitter)

    def init_databases(self):
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn")
        db_ms.setHostName(MYSQL_HOST)
        db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER)
        db_ms.setPassword(MYSQL_PASS)

        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST)
        db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER)
        db_pg.setPassword(PG_PASS)

        if not db_ms.open() or not db_pg.open():
            QMessageBox.critical(self, "Ошибка", "Не удалось подключиться к базе данных.")
            return False
        return True


if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setFont(QFont("Segoe UI", 10))
    app.setStyle("Fusion")
    window = MainWindow()
    window.show()
    sys.exit(app.exec())