# app.py

from flask import Flask, request, render_template_string
import csv

def create_app():
    app = Flask(__name__)

    # HTML template for the form
    HTML_TEMPLATE = '''
    <!DOCTYPE html>
    <html>
    <head>
        <title>Alcohol Consumption Checker</title>
    </head>
    <body>
        {% if user_type == 'admin' %}
            <h2>Admin: Add User Details</h2>
            <form method="POST">
                Name: <input type="text" name="name"><br>
                Age: <input type="number" name="age"><br>
                Gender: <select name="gender">
                            <option value="male">Male</option>
                            <option value="female">Female</option>
                            <option value="other">Other</option>
                        </select><br>
                <input type="submit" value="Submit">
            </form>
        {% else %}
            <h2>Check Alcohol Consumption Eligibility</h2>
            <form method="POST">
                Name: <input type="text" name="name"><br>
                <input type="submit" value="Check">
            </form>
        {% endif %}
        {% if result is not none %}
            <h3>{{ result }}</h3>
        {% endif %}
    </body>
    </html>
    '''

    @app.route('/run-tests', methods=['GET', 'POST'])
    def run_tests_view():
      if request.method == 'POST':
        import test_app
        test_results = test_app.run_tests()
        return render_template_string('''
            <h1>Test Results</h1>
            <pre>{{ test_results }}</pre>
            <a href="/run-tests">Run Again</a>
        ''', test_results=test_results)
      return render_template_string('''
        <h1>Run Tests</h1>
        <form method="post">
            <input type="submit" value="Run Tests">
        </form>
       ''')

    @app.route('/<user_type>', methods=['GET', 'POST'])
    def index(user_type):
        if request.method == 'POST':
            name = request.form['name']

            if user_type == 'admin':
                # Admin adding a new user to the database
                age = int(request.form['age'])
                gender = request.form['gender']
                with open('users.csv', mode='a', newline='') as file:
                    writer = csv.writer(file)
                    writer.writerow([name, age, gender])
                result = "User added to database."
            else:
                # Non-admin checking for alcohol consumption eligibility
                with open('users.csv', mode='r') as file:
                    reader = csv.reader(file)
                    for row in reader:
                        if row[0].lower() == name.lower():
                            result = f"{name}, you are {'allowed' if int(row[1]) >= 21 else 'not allowed'} to consume alcohol."
                            break
                    else:
                        result = "User not found in database."
            return render_template_string(HTML_TEMPLATE, user_type=user_type, result=result)

        return render_template_string(HTML_TEMPLATE, user_type=user_type, result=None)

    return app

if __name__ == '__main__':
    app = create_app()
    app.run(debug=True)

