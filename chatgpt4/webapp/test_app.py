# test_app.py

import unittest
import os
import csv
import io
from contextlib import redirect_stdout
from app import create_app

class FlaskTestCase(unittest.TestCase):
    def setUp(self):
        self.app = create_app()
        self.client = self.app.test_client()
        self.app.config['TESTING'] = True

        # Setup a clean database file
        with open('users.csv', mode='w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['Name', 'Age', 'Gender'])

    def add_and_print_user(self, name, age, gender):
        print(f"Adding User: Name: {name}, Age: {age}, Gender: {gender}")
        response = self.client.post('/admin', data={
            'name': name,
            'age': age,
            'gender': gender
        })
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'User added to database.', response.data)
        print("Result: User added successfully.\n")

    def check_user_eligibility(self, name):
        print(f"Checking Eligibility for: Name: {name}")
        response = self.client.post('/user', data={'name': name})
        self.assertEqual(response.status_code, 200)
        if b'not allowed' in response.data:
            print(f"Result: {name} is NOT allowed to consume alcohol.\n")
        else:
            print(f"Result: {name} is allowed to consume alcohol.\n")

    def test_user_operations(self):
        # Admin adding users
        self.add_and_print_user('John Doe', 25, 'male')
        self.add_and_print_user('Jane Doe', 19, 'female')

        # Non-admin checking eligibility
        self.check_user_eligibility('Jane Doe')

    def tearDown(self):
        # Clean up after each test case
        os.remove('users.csv')

def run_tests():
    buffer = io.StringIO()
    with redirect_stdout(buffer):
        tests = unittest.TestLoader().loadTestsFromTestCase(FlaskTestCase)
        unittest.TextTestRunner(stream=buffer, verbosity=2).run(tests)
    return buffer.getvalue()

if __name__ == '__main__':
    run_tests()

