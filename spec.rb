require 'rspec'
require 'open3'
require 'tmpdir'

def blogctl(args)
  go_args = "ci-build/blogctl.test -test.run '^Test_'"\
    " -test.coverprofile=coverage/#{rand 0..10_000}.out"

  Open3.capture3({'TEST_ARGS' =>  args}, go_args)
end

describe 'CLI' do

  context 'Trivial test cases' do
    it 'shows a version' do
      out, err, _ = blogctl("--version")
      expect(err).to eq ""
      expect(out).to match(/BuildVersion:/)
      expect(out).to match(/BuildDate:/)
    end
  end

  context 'Basic test cases' do

    let (:blog_path) { "/tmp/blog#{rand 0..10_0000}" }
    
    after(:each) do
      FileUtils.rm_rf(blog_path)
    end

    it 'initializes a new blog' do
      _, err, _ = blogctl("init")
      expect(err).to eq ""
      expect(File.exist?("blog")).to be_truthy
    ensure
      FileUtils.rm_rf("blog")
    end

    it 'initializes a new blog with custom path' do
      _, err, _ = blogctl("init -p #{blog_path}")
      expect(err).to eq ""
      expect(File.exist?(blog_path)).to be_truthy
    end

    it 'fails to initializes a blog twice' do
      _, _, status = blogctl("init -p #{blog_path}")
      expect(status.success?).to be_truthy
      expect(File.exist?(blog_path)).to be_truthy
      
      _, err, status = blogctl("init -p #{blog_path}")
      expect(err).to match(/Blog environment already exists/)
      expect(status.success?).to be_falsy
    end
  end
end
